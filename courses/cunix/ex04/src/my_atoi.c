/* converts str ==> int */

int my_atoi(const char *nptr)
{
    int i = (nptr[0] == '-') ? 1 : 0;
    int num = 0;
    for (; nptr[i] >= '0' && nptr[i] <= '9'; i++)
    {
        num = 10 * num + (nptr[i] - '0');
    }
    return (nptr[0] == '-') ? -num : num;
}
