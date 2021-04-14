The aim of this project is to provide a playground for recommender algorithms.
You are totally free how to use it.
- use the provided client-server to extract rules from your data
- use the code for experiments or educational reasons
- extend the code-base or use it as basis for your own system

One thing should be mentioned though: The code is not meant to be production ready
nor does the author take any responsibility for mistakes or bugs.

# Introduction

One of the most often used applications of machine learning arises in
recommender systems.
Their goal is to extract rules from given sets of items which lets predict
the 'purchase' of an item based on items already in the 'basket'.
That is, a user having for instance the items {1, 2, 3} in the basket, how probable is
it the user might be interested in item 7 as well.
The user might not even know of item 7 and by automatically suggesting it the system kind of
generates its own data-driven marketing strategy.
Usage of this not only is restricted on increasing revenue but also can be
used to suggest research articles or books a user might be interested in.
In other words, it can mitigate the effort to search for specific subjects.

## Implementation

The server is written in Go (https://golang.org/) and uses the Gin-framework together with
a websocket implementation of Gorilla (https://github.com/gorilla/websocket).
In addition it connects to a MySql database (https://www.mysql.com/) via Gorm (https://gorm.io/) in order
to store extracted rules.

The client is using React (https://reactjs.org/), Redux (https://redux.js.org/), React-Router (https://reactrouter.com/)
and Typescript (https://www.typescriptlang.org/).
The layout and styles is based on Antd (https://ant.design/) and the project has been initialized
by use of create-react-app (https://create-react-app.dev/).

## Running the App

There are two ways of running the app:
1. docker-compose
- put the docker-compose.yml into a folder or just clone the project
- ensure to have docker-compose installed on your system
- in a terminal type: docker-compose up

This downloads and starts the app image as well as a mysql database.
After this it automatically runs the app and you can access it at:
http://localhost:8080/app/data-import

2. build the project
- this require to have a mysql database listing on port 3306 configured with a database 'recommender'
  and accessible via: user: root, pwd: secret
- you will need to have yarn (https://yarnpkg.com/), nodejs (https://nodejs.org/en/) and Go (https://golang.org/)
  installed on your system
- clone the project and from within this target folder do the following comments in a terminal: <br>
  go build -o app . <br>
  cd client  <br>
  yarn build  <br>
  cd ..  <br>
  ./app  <br>
- this first builds the backend, then the client and then starts the app
- from this moment on, it tries to connect to the database

## Usage

The UI is trying to be quite self-explanatory.
The basic idea is:
- you drop your data, which require a specific format, into the system
- you select a model and some parameter
- you hit 'fit model', which starts to build the model and shows progress information
- you go to 'model statistics', which gives you some very basic view onto the outcome
- you go to 'model usage', which lets you test the model by producing recommendations

Note: Since especially the 'apriori' algorithm almost has exponential complexity, do not use
it with very large sets of data. In worst case, you will receive an out-of-memory.

## Example Data

The folder 'test-data' contains a correct formatted file with sample data: input_data.csv
You can just drop them into the UI, choose 'apriori' as model and a support of 3.
Alternative you can choose 'cosine' and a support of 2.
Note, these data actually present a bad case since not many items do appear in different orders.
The same folder contains a small python script which can help to transform data into the required format.


The following intends to give some easy introduction into both algorithms which are supported
by the app.

## Apriori Algorithm

The basic idea is best explained by an example.
Assume a store has data of purchases which contain item-ids like, <br>
{  <br>
{1,2,4},  <br>
{1,2,5,6}, <br>
{1,5,6}, <br>
{6,7} <br> 
}  <br>

We suppose the items within an item-set are ordered ascending.
First extract items which are deemed 'frequent', that is occurring with a frequency
larger than a fixed value. This value is referred to as 'support'. For this example we choose a support of 2.

We start of extracting all different items into an ordered list:<br>
{1,2,4,5,6,7}

Next we compute the frequencies for each single item: <br>
1: 3  <br>
2: 2  <br>
4: 1 <br>
5: 2 <br>
6: 3 <br>
7: 1 <br>

We drop all items which are below the value of support. This leaves us with: <br>
{1,2,5,6}

This we actually already can interpret as first set of rules. That is, starting
with an empty list, above items have higher probability of being chosen compared to others.
Rules: <br>
[] -> {1} <br>
[] -> {2} <br>
[] -> {5} <br>
[] -> {6} <br>

Next, based on above items we construct all possible item lists of length 2: <br>
{1, 2} <br>
{1, 5} <br>
{1, 6} <br>
{2, 5} <br>
{2, 6} <br>
{5, 6} <br>
Note, how we exploit the ordering here.

Then we compute frequencies of these lists as they are contained as subset in the original data set: <br>
{1, 2}: 2 <br>
{1, 5}: 2 <br>
{1, 6}: 2 <br>
{2, 5}: 1 <br>
{2, 6}: 1 <br>
{5, 6}: 2 <br>

As before we remove all item-lists which are below support: <br>
{1, 2}: 2 <br>
{1, 5}: 2 <br>
{1, 6}: 2 <br>
{5, 6}: 2 <br>

We would be tempted to extract rules like, {1} -> 2 or {5} -> 1 from that. But in order
to ensure statistical relevance we have to observe its confidence. For the rule {1} -> 2 this
is computed by: 2 / 3
We have 2 occurrences of item '2' based on 3 occurrences of item '1' altogether.
So in 2/3 of all cases, having '1' implies having '2'.
One can interpret this as probability of a binary distribution and so we only are interested in
rules which contribute to a confidence larger than 1/2.
We compute for all rules their confidences: <br>
{1} -> 2: 2/3 <br>
{1} -> 5: 2/3 <br>
{1} -> 6: 2/3 <br>
{5} -> 6: 2/2 <br>
{2} -> 1: 2/2 <br>
{5} -> 1: 2/2 <br>
{6} -> 1: 2/3 <br>
{6} -> 5: 2/3 <br>

All rules have confidence over 1/2 and so we would recognize them as statistical relevant.

We proceed by constructing item list of length 3 the same way we have done for length 2. In detail,
we combine all single frequent items, {1,2,5,6}, with item-list of length 2 which surpass the support: <br>
{1, 2, 3}: 0  <br>
{1, 2, 5}: 1 <br>
{1, 2, 6}: 1 <br>
{1, 5, 6}: 2 <br>

Here we have computed the item-set's frequencies at the same time.
As we see, only {1, 5, 6} is surpassing the support with value 2.
We state the rules with their confidence: <br>
{1, 5} -> 6: 2/2 <br>
{1, 6} -> 5: 2/2 <br>
{5, 6} -> 1: 2/2 <br>

All rules have confidence larger than 1/2 and thus we would add them.

Here the algorithm stops, since from {1, 5, 6} the only item set of length 4 we can construct is
{1, 5, 6, 7}. But this has frequency 0.

You might ask why we only have to consider {1, 5, 6, 7} and not for instance {1, 3, 6, 7}.
The reason for this lies in that we constructed item sets in lexicographical order:
Assume {1, 5, 6, 7} would occur as well with a frequency larger or equal to the support.
This would imply {1, 5, 6} as a subset would have a frequency of at least the same size.
But then it had been contained in the list of item set with length 3 which it didn't.
Or in other words, any item set of length 3 which lies lexicographical before {1, 5, 6} has been verified
with respect to its frequency and hence we do not need to re-consider when generating item sets of length 4.

Our algorithm altogether has extracted the rules: <br>
[] -> {1} <br>
[] -> {2} <br>
[] -> {5} <br>
[] -> {6} <br>
{1} -> 2 <br>
{1} -> 5 <br>
{1} -> 6 <br>
{5} -> 6 <br>
{2} -> 1 <br>
{5} -> 1 <br>
{6} -> 1 <br>
{6} -> 5 <br>
{1, 5} -> 6: 2/2 <br>
{1, 6} -> 5: 2/2 <br>
{5, 6} -> 1: 2/2 <br>

Although this algorithm might appear as the most natural and correct way of extracting rules it has
a sharp drawback: its complexity. First of all in order to compute item-set frequencies one has to
iterate the entire data. This iteration can occur quite often especially when one is facing too many
similar frequent item sets. In theory complexity can become almost exponential in the number of items.


## Cosine Algorithm

Let us consider the so call item-matrix:
{ <br>
  1, 1, 0, 1, 0, 0, 0 <br>
  1, 1, 0, 0, 1, 1, 0 <br>
  1, 0, 0, 0, 1, 1, 0 <br>
  0, 0, 0, 0, 0, 1, 1 <br>
} <br>

Each row presents a purchase and each column an item.
So the first row {1, 1, 0, 1, 0, 0, 0} would state:
Items 1, 2, 4 have been purchased and the others not.

The cosine-algorithm computes for each pair of item-columns the cosine.
This is defined by:
(a,b)/(|a||b|)

Here (a,b) is the euclidean inner product of the columns a and b, whereas
|a|, |b| denotes the euclidean norms of a and b respectively.

Our implementation deviates slightly from this by computing the pseudo-cosine.
That is, instead of using the euclidean norm in the dominator, we take the number
of purchases in which either item a or item b is involved.

As an example let us consider the first and second item column:
The inner product is given by: 1 * 1 + 1 * 1 + 1 * 0 + 0 * 0 = 2
The number of purchases is given by: 3
This makes the pseudo-cosine having the value: 2/3
In other words, out of 3 purchases, 2 have contained both items.

Given a set of items in the basket, recommendations are produced
by looking up for these items all other items which have a large
pseudo-cosine when being paired with one of these items.

The advantage of this algorithm compared to apriori is the ability to
run on large (if not huge) data sets.
The complexity is quadratic in the number of different items and linear
in the number of purchases.
So, if you have a huge number of purchases but only a few thousand of items,
this provides a perfect candidate.




















