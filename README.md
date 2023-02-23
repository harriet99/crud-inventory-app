# crud-inventory-app

### This application is a simple CRUD app built with Go and Gin. This lets you track your inventory, add new items, remove items, change items info and export to CSV.

 ### To run the application with docker:
 ```shell
 docker build -t inventory-app .
 ```
 ```shell
 docker run -p 8080:8080 inventory-app
 ```
 Then the code will run in http://localhost:8080

### Endpoints:
- (POST) localhost:8080/item
  - Post with a json body containing an id, name, quantity and unit_price. Returns the item added with a 201 code.
- (GET) localhost:8080/item
  - Returns a json with the list of the items and a 200 code.
- (GET) localhost:8080/item/1
  - Returns a json with the item wanted, and a 200 code.
- (DELETE) localhost:8080/item/1
  - Delete an item by its ID. Returns 200 if item was removed.
- (PATCH) localhost:8080/item/2
  - Body needs to contain a json with name, quantity and unit_price. The info of this item will be changed. Returns 200.
- (GET) localhost:8080/item/csv
  - Returns a CSV file.

### Database
- In this code, SQLite is used as the database for storing and retrieving inventory items. The advantage of using SQLite is its simplicity, ease of use, and low overhead. It requires no setup or configuration, and it can handle most small-scale data storage and retrieval tasks efficiently. Additionally, it is open-source and has a large user community, which means that there is plenty of support available.
