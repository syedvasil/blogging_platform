## Prepopulating some data

- Run mongosh
Step 1
```
$ use blogging_platform
```

Step 2 creating user including and Admin one
```
$ db.users.insertMany([ { _id: ObjectId("60c72b2f4f1a4b5d1c8f3b4d"), username: "JohnDoe", password: "password1", created_at: new Date(), role:"user" },
{ _id: ObjectId("60c72b2f4f1a4b5d1c8f3b4e"),role:"user", username: "JaneDoe", password: "password2", created_at: new Date() },
{ _id: ObjectId("60c72b2f4f1a4b5d1c8f34be"),role:"admin", username: "admin", password: "password3", created_at: new Date() }] )
```

verify that the users are created 
