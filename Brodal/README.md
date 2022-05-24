# Brodal heap

## My process
![alt text](https://xkcd.com/comics/good_code.png)

Image source: https://xkcd.com/844/

## Major pains during implementation

1. Figuring out `guide` data structure
   - this data structure is responsible for handling insertion/deletion of T1 and T2 root children in O(1)
   - it consists of 2 arrays and an integer T
   - in short, guide maintains an array of elements which must satisfy an invariant that $$
   x_{i} \le T, 1 \le i \le n
   $$, where n is the length of the array
   - for a better description of guides, read ![On the Power of Structural Violations in Priority](https://crpit.scem.westernsydney.edu.au/abstracts/CRPITV65Elmasry.html)
