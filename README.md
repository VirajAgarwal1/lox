# L O X 

## POSSIBLE FUTURE IMPROVEMENTS:

### lexer

1. **dfa**
    1. A Trie-based approach would have been faster for exact keywords 
    2. Some grouping would also help, like first group based on digits, letters or operators, then sub-group further. *[The 1st point in scanner can make this point obsolete]*
2. **scanner** 
    1. DFAs which start giving INVALID, do run them from the next run