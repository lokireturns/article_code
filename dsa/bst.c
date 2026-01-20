#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>

typedef struct Node node_t;

struct Node {
    uint32_t value;
    node_t* left;
    node_t* right;
};

node_t* init_bst(uint32_t root_val){
    node_t *new_tree = (node_t*)malloc(sizeof(node_t));
    new_tree->value = root_val;
    new_tree->left = NULL;
    new_tree->right = NULL;
    return new_tree;
}


node_t* add_node(node_t *root, node_t *node){
    
    if (root == NULL){
        root = node;
        return root;
    }   

    if(node->value < root->value){
        root->left = add_node(root->left, node);
    } else if (node->value > root->value){
        root->right = add_node(root->right, node);
    } else {
        perror("Could not add node");
        exit(EXIT_FAILURE);
    }
    
    return root;
}

void print_prompt(){ printf("Enter next node value \n:");};
void read_value(int * buffer){
    scanf("%d", buffer);
}

int main(int argc, char *argv[]){
    int root_val = 0;
    printf("Enter root node val: \n");
    scanf("%d", &root_val);
    node_t *tree = init_bst(root_val);
    int buffer;
    while (true)
    {
        print_prompt();
        read_value(&buffer);
        node_t *new_node = (node_t*)malloc(sizeof(struct Node));
        new_node->left = NULL;
        new_node->right = NULL;
        new_node->value = buffer;
        add_node(tree,new_node);
    }
    

    return 0;
};