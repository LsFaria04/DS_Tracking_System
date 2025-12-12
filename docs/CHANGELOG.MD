# CHANGELOG

## Sprint 3

In this sprint, the main priority was integrating the service with the remaining platform components and improving the overall testing strategy. This included Pub/Sub integration with both the notifications service and the checkout service, as well as aligning the service’s visual design with the overall website theme. Unit test coverage was also increased to ensure that every component behaves as expected. Additionally, a small feature was introduced that allows users to cancel their order while it is still in the processing phase.

Focusing first on the visual improvements, all components and pages of the service were updated to match the website’s overall design language. The changes can be seen in the images provided.

![sprint3](Images/Changelog/sprint%203/visuals_order_page.png)

![sprint3](Images/Changelog/sprint%203/visuals_order_carbon.png)

![sprint3](Images/Changelog/sprint%203/visuals_tracking_history.png)

Users can now cancel their orders while they are still being processed.

![sprint3](Images/Changelog/sprint%203/cancel.png)

To prevent accidental cancellations, a confirmation modal appears after the user clicks “Cancel” to ensure they truly intend to cancel the order.

![sprint3](Images/Changelog/sprint%203/cancel_modal.png)

As mentioned in the introduction, the service is now integrated with both the checkout and notifications services. However, these changes are not immediately visible in the user interface.


## Sprint 2

In this sprint, we introduced several new features to enhance the user experience on the orders page. First, the service is now accessible directly from the landing page once the user is authenticated.

![sprint2](Images/Changelog/sprint%202/landing_page.png)

When users click on “My Order” on the landing page, they are now directed to the new orders page, which displays all of their orders and provides options to sort by date or filter by status.

![sprint2](Images/Changelog/sprint%202/orders_page.png)

When accessing an order page, users can now view the real-time estimated delivery based on the order’s last known location, including potential delays caused by traffic.

![sprint2](Images/Changelog/sprint%202/estimated_delivery.png)

If users enter the wrong address during checkout, they can still update their delivery address while the order is in “Processing” by using the new order update feature.

![sprint2](Images/Changelog/sprint%202/order_update.png)

In line with MIPS’s sustainable image, the order page now includes a carbon footprint counter that displays emissions based on the means of transport used in the delivery.

![sprint2](Images/Changelog/sprint%202/carbon_footprint.png)

To improve transparency in the delivery process, the order page now displays the transaction hash for each order update, along with a link to the corresponding transaction on the Sepolia test network.

![sprint2](Images/Changelog/sprint%202/blockchain.png)

Other, less visible features introduced in this sprint include:
- Automatic order status updates via Pub/Sub. Currently, this is only tested with mock interactions, as no other team has integrated with the feature yet.
- New order reception via Pub/Sub. This functionality also relies on mock interactions, since the checkout team did not coordinate with us despite several attempts from our team.
- User notifications. Although implemented, this feature is not yet visible due to delays from the notifications team.

In addition, the team integrated with the broader MIPS infrastructure, which involved migrating to technologies such as React 18.0.0 and establishing communication with the main Pub/Sub system.



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
