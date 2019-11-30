# retail-sample

two main areas: in & out

`in` has sub areas:
* to configure item types (role: admin)
* to enter raw materials (stock item): type and quantity (role: user)

`out` has sub areas: 
* configure finished product (FP) components (role: admin)
* sell finished product (role: user)

## Details

a) item type:
- name - unique
- code - unique

actions:
* add if not present
* remove if not used
* disable if used

b) stock item: 
- code - points to a certain `enabled` `item type code`
- quantity - scalar value (integer) above zero and negative for correction/discard

action: can add new and if code is present then sum the quantity

c) finished  product type:
- name - unique
- a list of components, where each one has:
    * an `item type`
    * a quantity

d) finished product out
- type `finished product type`
- quantity - positive `int`

should update the stock - subtract the corresponding amounts