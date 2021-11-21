#include "test.h"

/* program to calc length of a string */

unsigned int my_strlen(char *str)
{
    unsigned int char_count = 0;
    unsigned int pos = 0;

    while (str[pos++] != '\0')
    {
        char_count++;
    }

    return char_count;
}
