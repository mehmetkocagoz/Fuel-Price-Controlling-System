# Brentoil Project Purpose
My primary objective is to establish a database encompassing Brent oil prices, diesel fuel prices, and the USD/TRY exchange rate. I will acquire this data by scraping information from various sources, including Bloomberght, TPDD.com.tr, and others. Subsequently, I intend to employ this dataset to predict fluctuations in fuel prices.

I have successfully achieved my initial goal. At this point, I am expanding my knowledge in data manipulation and data visualization. For the second phase of the project, I am considering the use of Python, as Python offers robust capabilities for data analysis, machine learning, and a wide range of scientific computing tasks.
## How to setup the project
Follow these steps to set up the project on your local machine:
### Clone the repository:
```

git clone https://github.com/mehmetkocagoz/brentoil.git

```
### Initialize go module
```
go mod init example.com/myproject

```
### Import dependencies
```
go mod tidy

```
### Update import parts
![importproblem](https://github.com/mehmetkocagoz/brentoil/assets/103457586/81a2c121-2177-4d32-8b73-d347dc271057)

Our directories is not same. So you need to change this parts. In the example go module above you should change your code like:
```
example.com/myproject/brentoil/database

```

## How to create your database
Follow these steps to set up the project on your local machine:

### Create a table
```postgresql
CREATE TABLE IF NOT EXISTS public.pricedata
(
    "timestamp" bigint,
    brentoilprice real DEFAULT 0,
    fuelprice real DEFAULT 0,
    date_column date DEFAULT (2222-22-02),
    exchange_column real DEFAULT 0
)
```

### Create a function
If you want to see timestamp's value in format (YYYY-MM-DD) you should use this function.
```
CREATE OR REPLACE TRIGGER update_date_column
    BEFORE INSERT
    ON public.pricedata
    FOR EACH ROW
    EXECUTE FUNCTION public.update_date();
```

### Create a trigger
It triggers on every insert action and update date_column
```
CREATE OR REPLACE TRIGGER update_date_column
    BEFORE INSERT
    ON public.pricedata
    FOR EACH ROW
    EXECUTE FUNCTION public.update_date();
```

## How to use cloned code to populate the database

### Connect to your database
In database.go file

![const](https://github.com/mehmetkocagoz/brentoil/assets/103457586/a9c9d61a-15e2-4a45-b157-3cad23fb6835)

As you can see there is a set of constants for my database.
You must change these values to connect your database.
Except changing values, you don't have to do anything to make connection.
### Use databaseFiller() function
In main.go call databaseFiller() function.

### Keep the database updated
In main.go whenever you call databaseUpdater() function. If new data comes from API, will be inserted to database.
