import pandas as pd

df = pd.read_csv("test-data/order_products__prior.csv", delimiter=",") # nrows=1000
df = df[["order_id", "product_id"]] \
      .groupby(["order_id"]) \
      .aggregate({"product_id": lambda x: ",".join(x.values.astype(str))})
print(df.head())

df.iloc[0:1000].to_csv("test-data/input_data.csv", mode="w", index=False, header=False, sep=" ")
