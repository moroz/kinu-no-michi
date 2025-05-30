begin;
  truncate products cascade;

\copy products (id, title, slug, base_price_eur, description, image_url) from stdin csv header;
id,title,slug,base_price_eur,description,image_url
019709a2-5c37-73e2-a05b-9ee9f8a470b5,Ocean Blue Shirt,ocean-blue-shirt,50,Ocean blue cotton shirt with a narrow collar and buttons down the front and long sleeves. Comfortable fit and tiled kalidoscope patterns. ,https://burst.shopifycdn.com/photos/young-man-in-bright-fashion_925x.jpg
01971efc-d170-7664-9ae8-82f386ff59fe,Classic Varsity Top,classic-varsity-top,60,"Womens casual varsity top, This grey and black buttoned top is a sport-inspired piece complete with an embroidered letter. ",https://burst.shopifycdn.com/photos/casual-fashion-woman_925x.jpg
\.

commit;
