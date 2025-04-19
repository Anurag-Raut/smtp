CREATE TABLE IF NOT EXISTS emails(
  id SERIAL PRIMARY KEY,
  "sentAt" DATE,
  "from" TEXT REFERENCES users(id),
  "to" TEXT ,
  data text
  
) 

CREATE INDEX idx_emails_from ON emails("from");
