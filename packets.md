# Update from server
|| byte delimiter
. empty (doesn't matter)
just for readability

## game params
|| grid_width || grid_height || tot_bombs || time || player ||
Tutto 2 byte tranne player
Player 1 byte: 
- 0: lobby owner
- 1: lobby guest

## updates
|| 0 0 game_over won .... || x || y || nnnn f ... ||

x = 2 bytes
y = 2 bytes
nnnn = cell label
f = Has next

last 5 bytes repeat until "has next" is false

## flags
|| 0 1 f ..... || x || y ||

x = 2 bytes
y = 2 bytes
f = flag value

# From client
|| x || y || ....... f ||
x, y: two bytes each
f: the move is a flag
