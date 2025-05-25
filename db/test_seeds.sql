begin;
  truncate products cascade;

\copy products (id, title, slug, base_price_eur, description, image_url) from stdin csv header;
id,title,slug,base_price_eur,description,image_url
019709a2-5c37-73e2-a05b-9ee9f8a470b5,Ocean Blue Shirt,ocean-blue-shirt,50,Ocean blue cotton shirt with a narrow collar and buttons down the front and long sleeves. Comfortable fit and tiled kalidoscope patterns. ,https://burst.shopifycdn.com/photos/young-man-in-bright-fashion_925x.jpg
\.

commit;
