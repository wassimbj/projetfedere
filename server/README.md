# Beamaan backend

## Stack

- **Golang** Backend language
- **Postgres** for storing data (users, transactions, ...)
- **Redis** as a queue broker and caching

## Folder Structure

```txt
   server/
   ├── api
   ├────── handlers => api routes handlers
   ├── config => app config
   ├── core => important code for the app
   ├── db => database connection
   ├── docker => dockerfiles
   ├── queue => queue system (using machinery)
   ├────── handlers => queue tasks handlers
   ├── services => DB queries functions
   ├── utils => utility functions 
```
