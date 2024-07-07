#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MAXVAL 255
#define STEP 17

static void print_brightness_level(int value)
{
    int len, filled;

    len = MAXVAL / STEP;
    filled = value / STEP;

    printf("%2d/%2d ", filled, len);

    for (int i = 0; i < len; i++) {
        char c = '-';

        if (i < filled) {
            c = '+';
        }
        printf("%c", c);
    }

    printf("\n");
}

static void usage(FILE *stream)
{
    fprintf(stream,
            "Usage: backlight-control <COMMAND>\n"
            "Changes brightness of a screen\n"
            "\n"
            "COMMAND:\n"
            "    set <VALUE> sets brightness level\n"
            "    inc         increases brightness\n"
            "    dec         decreases brightness\n"
            "    show        shows brightness level\n"
            "\n"
            "VALUE:\n"
            "    integer (0-%d)\n", MAXVAL / STEP);
}

int main(int argc, char **argv)
{
    FILE *f;
    int value;

    argv++;
    argc--;
    if (argc < 1) {
        fprintf(stderr, "Not enough arguments\n\n");
        usage(stderr);
        return 1;
    }

    if (strcmp(*argv, "--help") == 0) {
        usage(stdout);
        return 0;
    }

    f = fopen("/sys/class/backlight/amdgpu_bl1/brightness", "r+");
    if (f == NULL) {
        perror("Failed to open brightness file");
        return 1;
    }

    if (fscanf(f, "%d", &value) != 1) {
        fprintf(stderr, "Failed to read brightness value\n");
        return 1;
    }

    rewind(f);

    if (strcmp(*argv, "inc") == 0) {
        value += STEP;
        if (value > MAXVAL) {
            value = MAXVAL;
        }
        fprintf(f, "%d", value);
        print_brightness_level(value);
        return 0;
    }

    if (strcmp(*argv, "dec") == 0) {
        value -= STEP;
        if (value < 0) {
            value = 0;
        }
        fprintf(f, "%d", value);
        print_brightness_level(value);
        return 0;
    }

    if (strcmp(*argv, "show") == 0) {
        print_brightness_level(value);
        return 0;
    }

    if (strcmp(*argv, "set") == 0) {
        argv++;
        argc--;
        if (argc < 1) {
            fprintf(stderr, "Not enough arguments\n\n");
            usage(stderr);
            return 1;
        }
        value = atoi(*argv) * STEP;
        if (value > MAXVAL) {
            value = MAXVAL;
        } else if (value < 0) {
            value = 0;
        }
        fprintf(f, "%d", value);
        print_brightness_level(value);
        return 0;
    }

    fprintf(stderr, "Unknown command: %s\n\n", *argv);
    usage(stderr);
    return 1;
}
