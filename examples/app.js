const { MongoClient } = require("mongodb");

const url = process.env.MONGO_URL;

async function ping(client) {
  try {
    const adminDb = await client.db("admin");
    const { ok } = await adminDb.command({ ping: 1 });
    console.log("Ping from NodeJS", ok ? "ok" : "ko");
  } catch (reason) {
    showPingError(reason);
  } finally {
    await client.close();
  }
}

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
  console.error(details);
};

async function run() {
  const client = new MongoClient(url, {
    connectTimeoutMS: 100,
    serverSelectionTimeoutMS: 100,
  });

  try {
    await client.connect();
    await ping(client); // Only attempt ping after connection succeeded.
  } catch (reason) {
    console.log("Could not connect to MongoDB:", reason);
    process.exit(1);
  } finally {
    // No need to close here: we closed in the ping.
  }
}

run().catch(console.error);
