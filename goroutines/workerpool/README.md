# Worker pool

What is a workerpool? It is a controlled way to run jobs among wokers. It normally has a fixed size of pool of worders.
So when there are too many jobs to do, the job queue piles up, but there is no dangeous of running out of resource.
If it is required, a pool can change the size dynamically. Minimal size of pool provides enough capacity with reserves to quickly process jobs. There should have a upper limit. This minimal size is chosen for the normal job load pattern.

For a dynamic worker pool to work, the queue should have enough capacity to hold job requests. A simple channel does not work. Channel, even buffered channel is not a storage but a way of communication.

Job request --> job_request_channel --> job_request_stack(FIFO) --> job dispatcher --> worker channel --> worker pool.

## Example one
A job dispatch queue has 5 jobs to dispatch. A size of two worker pool to process them.