# Go Concurrent Programming Challenge: Producer-Consumer with Multiple Buffers

## Description:

Imagine a production and consumption system, like a factory where products are manufactured (produced) and sent to a warehouse (buffer). Several workers (goroutines) produce items and place them in different buffers. Other workers remove items from the buffers and process them.

## Requirements:

Multiple Buffers: Implement a system with N buffers.

Producers: Create M goroutines as producers. Each producer should:
- Randomly choose a buffer to place the item.
- Produce items randomly (random production time).

Consumers: Create K goroutines as consumers. Each consumer should:
- Randomly choose a buffer to remove the item.
- Consume items randomly (random consumption time).

Synchronization: Use channels to synchronize production and consumption between goroutines.

Buffer Capacity: Each buffer has a maximum capacity of items.

Termination Condition: The program should stop after a certain amount of time or when all producers have finished producing.