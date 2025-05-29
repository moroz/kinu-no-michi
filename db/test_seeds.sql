begin;
  truncate products cascade;

\copy products (id, title, slug, base_price_eur, description, image_url) from stdin csv header;
id,title,slug,base_price_eur,description,image_url
019709a2-5c37-73e2-a05b-9ee9f8a470b5,Ocean Blue Shirt,ocean-blue-shirt,50,Ocean blue cotton shirt with a narrow collar and buttons down the front and long sleeves. Comfortable fit and tiled kalidoscope patterns. ,https://burst.shopifycdn.com/photos/young-man-in-bright-fashion_925x.jpg
019709a2-5c3a-7c63-b811-090eaedf0835,Classic Varsity Top,classic-varsity-top,60,"Womens casual varsity top, This grey and black buttoned top is a sport-inspired piece complete with an embroidered letter. ",https://burst.shopifycdn.com/photos/casual-fashion-woman_925x.jpg
019709a2-5c3a-75d4-97a5-64b49a954cab,Yellow Wool Jumper,yellow-wool-jumper,80,Knitted jumper in a soft wool blend with low dropped shoulders and wide sleeves and think cuffs. Perfect for keeping warm during Fall. ,https://burst.shopifycdn.com/photos/autumn-photographer-taking-picture_925x.jpg
\.

commit;
