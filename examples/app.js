const { MongoClient } = require("mongodb");

const url = process.env.MONGO_URL;

const ping = async (db) => {
  let adminDb = await db.db("admin");
  const { ok } = await adminDb.command({ ping: 1 }).catch(showPingError);
  console.log("Ping from NodeJS", ok ? "ok" : "ko");
  db.close();
};

const showPingError = async (reason) => {
  if (reason.constructor.name !== "MongoError") {
    throw reason;
  }

  // Just like the mongo shell.
  const details = {
    errmsg: reason.errmsg,
    code: reason.code,
    codeName: reason.codeName,
  };
  console.log(details);
};

async function run() {
  const client = new MongoClient(url, {
    useUnifiedTopology: true,
    connectTimeoutMS: 100,
    serverSelectionTimeoutMS: 100,
  });
  try {
    await client.connect().catch((reason) => {
      console.log("could not connect", reason);
      process.exit(1);
    });
    await ping(client);
  } finally {
    await client.close();
  }
}

run().catch(console.dir);
