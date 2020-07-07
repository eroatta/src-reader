db.createUser(
    {
        user: "reader_rw",
        pwd: "reader_rw_password",
        roles: [
            {
                role: "readWrite",
                db: "reader"
            }
        ]
    }
)
