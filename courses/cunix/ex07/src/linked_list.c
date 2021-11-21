#include <stdio.h>
#include <stdlib.h>

typedef struct node
{
    void *data;
    struct node *next;
} node_t;

node_t *list_create(void *data);

void list_destroy(node_t **head, void (*fp)(void *data));
void list_push(node_t *head, void *data);
void list_unshift(node_t **head, void *data);

void *list_pop(node_t **head);
void *list_shift(node_t **head);
void *list_remove(node_t **head, int pos);

void list_print(node_t *head);
void list_visitor(node_t *head, void (*fp)(void *data));

node_t *list_create(void *data)
{
    node_t *root = (node_t *)malloc(sizeof(node_t));
    root->next = NULL;
    root->data = data;
    return root;
}

void list_destroy(node_t **head, void (*fp)(void *data))
{
    if (*head == NULL)
        return;

    node_t *ptr = *head;
    node_t *next;

    while (ptr != NULL)
    {
        fp(ptr->data);
        next = ptr->next;
        free(ptr);
        ptr = next;

    }

    *head = NULL;
}

void list_push(node_t *head, void *data)
{
    if (head == NULL){
        head = list_create(data);
        return;
    }
    node_t *ptr = head;

    while (ptr->next != NULL)
    {
        ptr = ptr->next;
    }

    ptr->next = list_create(data);
}

void list_unshift(node_t **head, void *data)
{
    if (*head == NULL){
        *head = list_create(data);
        return;
    }

    node_t *new_node = list_create(data);

    new_node->next = *head;
    *head = new_node;
}

void list_print(node_t *head)
{
    if (head == NULL)
        return;

    node_t *ptr = head;

    while (ptr != NULL)
    {
        printf("%s\n", (char *)(ptr->data));
        ptr = ptr->next;
    }
}

void *list_pop(node_t **head){

    if (*head == NULL)
        return NULL;

    void* data;
    node_t *ptr = *head;

    if(ptr->next == NULL){
        data = ptr->data;

        free(ptr->data);
        free(ptr);

        *head = NULL;
        return data;
    }

    while(ptr->next->next != NULL)
        ptr = ptr->next;

    data = ptr->next->data;

    free(ptr->next->data);
    free(ptr->next);

    ptr->next = NULL;
    return data;
}

void *list_shift(node_t **head){

    if(*head == NULL)
        return NULL;

    node_t* tmp = *head;
    void* data = tmp->data;

    *head = tmp->next;

    free(tmp);
    return data;
}

void *list_remove(node_t **head, int pos)
{
    node_t *current = *head;
    while (pos-- > 0)
    {
        current = current->next;
    }

    node_t *removed_node = current->next;

    current->next = removed_node->next;

    free(removed_node->data);
    free(removed_node);

    return removed_node;
}


void list_visitor(node_t *head, void (*fp)(void *data))
{
    node_t *node = head;
    while (node)
    {
        (*fp)(node->data);
        node = node->next;
    }
}
