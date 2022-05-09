# Brodal heap

## My process
![alt text](https://xkcd.com/comics/good_code.png)

Image source: https://xkcd.com/844/

## Problems I encountered
### 1. Guide implementation

Brodal described a _guide_ implementation [here](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwiRj5-br9L3AhVQ_bsIHY9ADEYQFnoECAcQAQ&url=https%3A%2F%2Fcs.au.dk%2F~gerth%2Fpapers%2Fsoda96.pdf&usg=AOvVaw33D6m3_qIfJD7dnt8TL46Y). He stated that the `REDUCE(i)` operation decreases $x_{i}$ by at least 2 and increases $x_{i+1}$ by at most 1. The "at most 1" part is my problem. I don't know if it's ok to just ignore missing increase or what...
