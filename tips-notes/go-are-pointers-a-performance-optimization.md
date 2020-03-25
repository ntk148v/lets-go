# Go: Are pointers a performance optimization?

Source: https://medium.com/@vCabbage/go-are-pointers-a-performance-optimization-a95840d3ef85

**TL;DR**: Pointers are not inherently a performance optimization.

* A pointer is a memory address. To access the value being pointed at, the program must follow the address to the beginning of the value. This is referred to as "dereferencing".
* In many cases, a pointer is smaller than the value being pointed at:
    * If the argument is a scalar type (`bool`, `int`, `float`,...), it's going to be less than or equal to the size of a pointer.
    * If the argument is a compound type (struct...), it's likely the pointer is smaller.
* Pointers can negatively affect performance:
    * Dereferencing pointers isn't free. It's not a huge cost, but it can add up.
    * Sharing data via pointers will likely cause the data to be placed in the "heap". The heap is section of memory for data that lives longer than a single function call. There is overhead to adding data to the heap and heap data can be only be cleaned up by the garbage collector.

## Stack vs Heap

### Stack: Function-local memory

* A function is called -> Create its own section of stack to store local variables.
* Stack size - compile time.
* Function is called -> the next area of free mem in stack -> the function.
* Function returns -> release -> the area is available for the next function call.

### Heap: Area of shared data

* Function local variables "disappear" after the function returns. The returned values are copied into the stack of the calling function.
* Pointers are returned, the pointed-at data needs to be placed somewhere outside that stack -> **heap**.
* Data into  heap -> requires memory from the runtime.
* Not enough heap space -> runtime will have to ask for additional  memory from the OS.
* A value has been placed in  the heap its needs to stay there  until no functions have a pointer to it anymore -> Garbage collector's job.

## Conclusion

### Advantages

* Pointer can avoid copying memory.
* Pointers allow you to share data.If you want a function to be able to modify the the data you’re passing it, a pointer is appropriate.
* Pointer can also  be useful when you need to distinguish between a zero value and an unset value.

### Disadvantages

* Additional levels of  indirection and increased work for the garbage collector.
