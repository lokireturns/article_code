#include <stdlib.h>
#include <stdio.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <unistd.h>

int main(int argc, char *argv[]) {

    int fd = open("some_file.txt", O_RDONLY);
    if(fd == -1){
        perror("open error");
        exit(EXIT_FAILURE);
    }

    struct stat filestats;
    if(fstat(fd, &filestats) == -1 ) {
        perror("stats");
        exit(EXIT_FAILURE);
    }

    // why are addresses returned as void in C?
    void *addr = mmap(NULL, filestats.st_size, PROT_READ, MAP_PRIVATE, fd, 0);
    if (addr == MAP_FAILED) {
        perror("mmap attempt");
        exit(EXIT_FAILURE);
    }

    // Now we can read the contents of the file but what about demand paging?
    // should we mmap to a smaller region of memory to see deman paging in action?
    printf("Contents of the file: %s\n", (char *)addr); // Dereference to array of chars?

    // We can safely close the file descriptor - mapping is preserved
    close(fd);

    // We can optionally unmmap
    if(munmap(addr, filestats.st_size) == -1) {
        perror("unmmap");
        exit(EXIT_FAILURE);
    }

    return 0;

};