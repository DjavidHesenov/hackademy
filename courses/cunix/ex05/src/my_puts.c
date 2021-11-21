#include <unistd.h>
/* output str (stdout) */
int my_puts(const char *s)
{
    char new_l = '\n';
    for (int i = 0; s[i] != '\0'; i++)
    {
        write(1, &s[i], sizeof(char));
    }
    write(1, &new_l, sizeof(char));

    return 1;
}
