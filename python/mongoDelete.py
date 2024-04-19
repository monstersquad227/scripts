from pymongo import MongoClient
from datetime import datetime


def batch_delete_documents(database_name, collection_name, batch_size):
    # MongoDB 连接字符串
    uri = "mongodb://192.168.1.35:27017"

    try:
        # 连接到 MongoDB
        client = MongoClient(uri)
        db = client[database_name]
        collection = db[collection_name]

        deleted_count = 0

        while True:
            time_threshold = datetime(2022, 12, 31)  # 替换为你希望的时间
            condition = {"CreateTime": {"$lt": time_threshold}, "LogLevel": 1} # mongo删除条件

            result = collection.delete_many(condition)
            count = result.deleted_count
            deleted_count += count
            
            if count < batch_size:
                break

        print(f"Total {deleted_count} documents deleted.")

    except Exception as e:
        print(f"An error occurred: {str(e)}")
    finally:
        if client:
            client.close()

if __name__ == "__main__":
    database_name = "Log"  # 替换为实际的数据库名称
    collection_name = "OperationLog"  # 替换为实际的集合名称
    batch_size = 1  # 批量删除大小（这里只是占位，实际删除大小由循环控制）

    batch_delete_documents(database_name, collection_name, batch_size)
