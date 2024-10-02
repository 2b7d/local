#define _POSIX_C_SOURCE 200112L

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/wait.h>

#define MAXVAL 255
#define STEP 17
#define STEPS (MAXVAL / STEP)

extern char **environ;
char *program_name;

void usage(FILE *stream)
{
    fprintf(stream,
            "Usage: %s <COMMAND>\n"
            "Change brightness of a screen\n"
            "\n"
            "COMMAND:\n"
            "    set <VALUE> set brightness level\n"
            "    inc         increase brightness\n"
            "    dec         decrease brightness\n"
            "    show        show brightness level\n"
            "VALUE:\n"
            "    integer (0 - %d)\n\n", program_name, STEPS);
}

void write_value(FILE *f, int value)
{
    int fd = fileno(f);

    if (fd < 0) {
        perror("fileno");
        exit(1);
    }

    if (ftruncate(fd, 0) < 0) {
        perror("ftruncate");
        exit(1);
    }

    fprintf(f, "%d", value);
}

int show_notification(int value)
{
    pid_t pid = fork();

    int wstatus;

    if (pid < 0) {
        perror("fork");
        exit(1);
    }

    if (pid == 0) {
        float max = MAXVAL;
        float fval = value;
        int perc = fval / max * 100;

        char progress[16];

        char *name = "/bin/dunstify";
        char *argv[] = {name, "-h", progress, "-r", "1", "-t", "500", "Brightness", 0};

        snprintf(progress, sizeof(progress), "int:value:%d", perc);

        if (execve(name, argv, environ) < 0) {
            perror("execve");
            exit(1);
        }
    }

    if (wait(&wstatus) < 0) {
        perror("wait");
        exit(1);
    }

    if (WIFEXITED(wstatus) != 1) {
        fprintf(stderr, "child did not exit normally\n");
        exit(1);
    }

    return WEXITSTATUS(wstatus);
}

int main(int argc, char **argv)
{
    int value;
    FILE *f;

    program_name = *argv;

    argc--;
    argv++;

    if (argc < 1) {
        usage(stderr);
        fprintf(stderr, "not enough arguments\n");
        return 1;
    }

    if (strcmp(*argv, "--help") == 0) {
        usage(stdout);
        return 0;
    }

    f = fopen("/sys/class/backlight/amdgpu_bl1/brightness", "r+");
    if (f == NULL) {
        perror("fopen");
        return 1;
    }

    if (fscanf(f, "%d", &value) != 1) {
        fprintf(stderr, "failed to read value from brightness file\n");
        return 1;
    }

    rewind(f);

    if (strcmp(*argv, "inc") == 0) {
        value += STEP;
        if (value > MAXVAL) {
            value = MAXVAL;
        }

        write_value(f, value);
        return show_notification(value);
    }

    if (strcmp(*argv, "dec") == 0) {
        value -= STEP;
        if (value < 0) {
            value = 0;
        }

        write_value(f, value);
        return show_notification(value);
    }

    if (strcmp(*argv, "set") == 0) {
        int new_value;

        argc--;
        argv++;

        if (argc < 1) {
            usage(stderr);
            fprintf(stderr, "not enough arguments\n");
            return 1;
        }

        new_value = atoi(*argv++);
        if (new_value > STEPS) {
            new_value = STEPS;
        } else if (new_value < 0) {
            new_value = 0;
        }

        new_value *= STEP;
        write_value(f, new_value);
        return show_notification(new_value);
    }

    if (strcmp(*argv, "show") == 0) {
        return show_notification(value);
    }

    usage(stderr);
    fprintf(stderr, "unknown command: %s\n", *argv);
    return 1;
}
