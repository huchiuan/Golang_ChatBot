db.createUser(
    {
        user: "edison",
        pwd: "edison",
        roles: [
            {
                role: "readWrite",
                db: "linemessage"
            }
        ]
    }
);
db.createCollection('linemessage');