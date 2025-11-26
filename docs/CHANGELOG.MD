# CHANGELOG

## Sprint 2
//TODO

## Sprint 1

In this sprint, we implemented most of the core functionalities of the order tracking system.

The first enhancement was the ability to view detailed order information, including the delivery address, estimated delivery date, and the list of products in the order.

![sprint1](Images/Changelog/sprint%201/order_info.png)

Secondly, we introduced a feature that allows users to verify the integrity of their orders using blockchain technology.

![sprint1](Images/Changelog/sprint%201/blockchain.png)

We also added a feature that allows users to view the position of their order on a map, along with its trajectory from the seller to the current location, including the intermediate stops it has passed through.

![sprint1](Images/Changelog/sprint%201/map.png)

Lastly, we introduced a history bar that displays all the key updates related to the order.

![sprint1](Images/Changelog/sprint%201/history.png)


Most of these features were already implemented or partially developed during Sprint 0. In Sprint 1, the team focused on ensuring that the Sprint 0 features were fully completed and properly integrated with the rest of the system. Unfortunately, integration with the checkout service was not possible due to delays from the team responsible for that component, even though our team had already partially configured Pub/Sub. On the positive side, we successfully integrated with the product service via the Jumpseller API. Additionally, the frontend code was migrated from Next.js to React to align with the rest of the frontend stack.
