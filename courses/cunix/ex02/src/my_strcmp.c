#include "test.h"

/* compare strings. <0 if str1<str2, 0 if str1=str2, >0 if str1>str2 */

int my_strcmp(char *str1, char *str2)
{
    for (; *str1 == *str2; str1++, str2++)
    {
        if (*str1 == '\0')
        {
            return 0;
        }
    }
    if (*str1 - *str2 < 0) 
    {	
        return -1;
    }
    else 
    {
        return 1;
    }
}
