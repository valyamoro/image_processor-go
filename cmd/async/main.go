package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    "image_processor-go/internal/core"
    "image_processor-go/internal/generator"
    "image_processor-go/internal/worker"
)

func main() {
    start := time.Now()
    
    imageCount := 100
    workerCount := 4
    inputDir := "../../test_input_async"
    outputDir := "../../test_output_async"
    targetFormat := core.FormatPNG
    
    fmt.Println("=== ASYNC MODE ===")
    fmt.Println("Generating test images...")
    if err := generator.GenerateTestImages(imageCount, inputDir); err != nil {
        fmt.Printf("Error generating images: %v\n", err)
        return
    }

    fmt.Printf("Starting worker pool (%d workers)...\n", workerCount)
    pool := worker.NewPool(workerCount, 1000)
    pool.Start()
    
    go monitorResults(pool.GetResultChan())
    
    statsTicker := time.NewTicker(2 * time.Second)
    go monitorStats(pool, statsTicker)
    
    fmt.Println("\nScanning files and creating jobs...")
    files, err := os.ReadDir(inputDir)
    if err != nil {
        fmt.Printf("Error reading dir: %v\n", err)
        return
    }
    
    for i, file := range files {
        if file.IsDir() {
            continue
        }
        
        inputPath := filepath.Join(inputDir, file.Name())
        outputName := fmt.Sprintf("img_%04d.%s", i+1, targetFormat)
        outputPath := filepath.Join(outputDir, outputName)
        
        info, err := os.Stat(inputPath)
        if err != nil {
            fmt.Printf("Error stating %s: %v\n", inputPath, err)
            continue
        }
        
        var priority core.Priority
        age := time.Since(info.ModTime())
        switch {
        case age > 365*24*time.Hour:
            priority = core.PriorityHigh
        case age > 30*24*time.Hour:
            priority = core.PriorityMid
        default:
            priority = core.PriorityLow
        }
        
        job := &core.Job{
            ID:          i + 1,
            Priority:    priority,
            InputPath:   inputPath,
            OutputPath:  outputPath,
            FileModTime: info.ModTime(),
            FileSize:    info.Size(),
            TargetFormat: targetFormat,
            Quality:     85,
        }
        
        pool.Submit(job)
        
        if (i+1)%50 == 0 {
            fmt.Printf("Submitted %d/%d jobs (queue: %d)\n", 
                i+1, len(files), pool.GetQueueLength())
        }
    }
    
    fmt.Printf("\nAll %d jobs submitted to queue\n", len(files))
    
    for pool.GetQueueLength() > 0 {
        time.Sleep(500 * time.Millisecond)
        fmt.Printf("Waiting... queue: %d\n", pool.GetQueueLength())
    }
    
    fmt.Println("\nStopping worker pool...")
    statsTicker.Stop()
    pool.Stop()
    
    elapsed := time.Since(start)
    printFinalStats(pool, elapsed, len(files))
}

func monitorResults(results <-chan string) {
    processed := 0
    for result := range results {
        processed++
        if processed%10 == 0 {
            fmt.Printf("[Results] %d: %s\n", processed, result)
        }
    }
}

func monitorStats(pool *worker.Pool, ticker *time.Ticker) {
    for range ticker.C {
        stats := pool.GetWorkerStats()
        busy := 0
        totalJobs := 0
        
        for _, s := range stats {
            if s.IsBusy {
                busy++
            }
            totalJobs += s.JobsDone
        }
        
        fmt.Printf("[Stats] Busy: %d/%d | Total jobs: %d | Queue: %d\n", 
            busy, len(stats), totalJobs, pool.GetQueueLength())
    }
}

func printFinalStats(pool *worker.Pool, elapsed time.Duration, totalFiles int) {
    stats := pool.GetWorkerStats()
    
    fmt.Printf("\n=== WORKER STATISTICS ===\n")
    totalJobs := 0
    for _, s := range stats {
        avgTime := time.Duration(0)
        if s.JobsDone > 0 {
            avgTime = s.TotalTime / time.Duration(s.JobsDone)
        }
        
        fmt.Printf("Worker #%d: jobs=%d, avg_time=%v, high=%d, mid=%d, low=%d\n",
            s.ID, s.JobsDone, avgTime, 
            s.HighPriorityJobs, s.MidPriorityJobs, s.LowPriorityJobs)
        totalJobs += s.JobsDone
    }
    
    fmt.Printf("\nTotal jobs processed: %d\n", totalJobs)

    fmt.Printf("\n=== FINAL RESULTS ===\n")
    fmt.Printf("Total time: %v\n", elapsed)
    fmt.Printf("Total files: %d\n", totalFiles)
    fmt.Printf("Speed: %.2f images/second\n", float64(totalFiles)/elapsed.Seconds())
}
