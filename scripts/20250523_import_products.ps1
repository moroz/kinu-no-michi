#!/usr/bin/env pwsh

$file = "./apparel.csv"
$dbURL = $Env:DATABASE_URL

if ($dbURL -eq "") {
  write-host "Error: DATABASE_URL is not set!"
  exit 1
}

# allocate a dynamic array of strings
$transformed = New-Object System.Collections.Generic.HashSet[string]

$lines = get-content $file
$lineCount = ($lines | Measure-Object -Line).Lines

for ($i = 1; $i -lt $lineCount; $i++) {
  $id = [Guid]::CreateVersion7()
  $transformed.Add("$id,$($lines[$i])") | out-null
}

$csv = $transformed -join "`n"

@"
begin;

create temporary table products_temp (
  id uuid not null primary key,
  handle text not null unique,
  title text not null,
  body text not null,
  tags text,
  price decimal not null,
  img_src text
);

\copy products_temp from stdin csv;
$csv
\.

select count(*) from products_temp;

delete from products;

insert into products (id, title, slug, base_price_eur, description, image_url)
select id, title, handle, price, body, img_src from products_temp;

commit;
"@ | psql $dbURL
