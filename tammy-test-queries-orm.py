import json
import os

from sqlalchemy import create_engine, insert, inspect, select
from sqlalchemy import Column, ForeignKey, Integer, MetaData, String, Table


rds_proxy_host = os.environ.get('RDS_PROXY_HOST')
db_name = os.environ.get('DB_NAME')
user_name = os.environ.get('USER_NAME')
pword = os.environ.get('PASSWORD')

db_engine = create_engine(f"postgresql://{user_name}:{pword}@{rds_proxy_host}")


class MyDB:
    def __init__(self):
        self.md = MetaData()
        self.author_table = Table(
            "author",
            self.md,
            Column("author_id", Integer, primary_key=True),
            Column("name", String, nullable=False),
        )
        self.book_table = Table(
            "book",
            self.md,
            Column('book_id', Integer, primary_key=True),
            Column('author_id', Integer, ForeignKey("author.author_id"), nullable=False),
            Column('title', String),
            Column('summary', String),
        )

    def select_author(self, author_name):
        rows = None
        with db_engine.connect() as conn:
            rows = conn.execute(
                select(self.author_table).where(
                    self.author_table.c.name == author_name
                )
            )
        return str(rows.first())

    def insert_author(self, author_name):
        with db_engine.connect() as conn:
            result = conn.execute(
                insert(self.author_table),
                [
                    {"name": author_name}
                ]
            )
            conn.commit()

    def select_book_by_title(self, book_title):
        rows = None
        with db_engine.connect() as conn:
            rows = conn.execute(
                select(self.book_table).where(
                    self.book_table.c.title == book_title
                )
            )
        return str(rows.first())

    def select_books_by_author(self, book_author_id):
        rows = None
        with db_engine.connect() as conn:
            rows = conn.execute(
                select(self.book_table).where(
                    self.book_table.c.author_id == book_author_id
                )
            )
        return str(rows.all())

    def insert_book(self, author_id, title, summary):
        with db_engine.connect() as conn:
            result = conn.execute(
                insert(self.book_table),
                [
                    {
                        "author_id": author_id,
                        "title": title,
                        "summary": summary,
                    }
                ]
            )
            conn.commit()


def lambda_handler(event, context):
    inspection = inspect(db_engine)
    table_names = inspection.get_table_names()

    my_db = MyDB()
    # my_db.insert_author("Isaac Asimov")
    # my_db.insert_book(2, "The Foundation", "Nerds are exiled to Terminus, what next?")
    author = my_db.select_author("Frank Herbert")
    books = my_db.select_books_by_author(1)

    return {
        'statusCode': 200,
        'table_names': table_names,
        'author': author,
        'books': books,
    }
