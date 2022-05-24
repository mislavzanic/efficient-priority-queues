# Brodal heap

## My process
![alt text](https://xkcd.com/comics/good_code.png)

Image source: https://xkcd.com/844/

## Major pains during implementation

$\sqrt{3x+1}$
1. Figuring out `guide` data structure
   - this data structure is responsible for handling insertion/deletion of $T_{1}$ and $T_{2}$ root children in $O(1)$
   - it consists of 2 arrays and an integer $T$
