# BPE model

This demonstrates how to train a tokenizer from scratch using BPE model.

It trains a tokenizer for Esperanto language from scratch using data from
`input` folder and save `vocab` and `merges` into `model` folder.

To run: 

```bash
# run training
go run . -mode=train

# run test
go run . -mode=test
```
