import configparser

import pyrfc
import csv
import os
import shutil
from datetime import datetime

from pyrfc import Connection

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


# распаковка ЕО
def main():
    global output
    logger = logging.getLogger("unpack_id")
    logger.setLevel(logging.INFO)

    ttime = datetime.now()
    # create the logging file handler
    fh = logging.FileHandler(
        "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/pyrfc_logging.log")
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    fh.setFormatter(formatter)

    # add handler to logger object
    logger.addHandler(fh)

    logger.info(f"Parsing cfg")
    config = configparser.ConfigParser()
    config.read(
        "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/sapnwrfc.cfg")
    config.sections()
    params_connection = config['connection']
    logger.info(f"Connecting to SAP RFC...")

    try:
        connection = pyrfc.Connection(**params_connection)
    #    result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
    #    print(result)
        work_order = None
        # work_order_name_f = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/work_order_name.csv"
        work_order_name_f = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_work_order_name.csv"
        # work_order_name_f = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_work_order_name.csv"

        dataArchive = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_archive/"
        # id_test_1000836.csv
        # file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/unpack_id.csv"
        # file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_unpack_id.csv"
        # file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_2_unpack_id.csv"
        file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_unpack_id.csv"
        # file_unpack_id = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_unpack_id.csv"

        # file_unpack_id_scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/unpack_id_scrap.csv"
        file_unpack_id_scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_unpack_id_scrap.csv"
        # file_unpack_id_scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_unpack_id_scrap.csv"

        # чтение списка выгруженных ЕО на редактирование
        rows = []
        with open(file_unpack_id, newline='') as file:
            csvreader = csv.reader(file, delimiter=',')
            header = next(csvreader)
            for row in csvreader:
                rows.append(row)
        # распаковка в цикле прочитанных ЕО
        for i in rows:
            output = connection.call('Z_IEXT_HU_UNPACKSNGLPOS', **{
                'UCODE': '21717',
                'PCODE': 'NEWPASSWORD1',
                'HUKEY': '0000000000' + i[0],  # '00000000000000015660',
                'ITEMUNPACK': {u'PACK_QTY': i[1]},

            })

        print(output)
        logger.info(f"Unpack ID: {output}")

        rowsScrap = []
        file_exist_id_scrap = os.path.exists(file_unpack_id_scrap)
        if file_exist_id_scrap:
            with open(file_unpack_id_scrap, newline='') as fileScrap:
                csvreader = csv.reader(fileScrap, delimiter=',')
                header = next(csvreader)
                for rowScrap in csvreader:
                    rowScrap.append(rowsScrap)
            for s in rowScrap:
                outputScrap = connection.call('Z_IEXT_HU_UNPACKSNGLPOS', **{
                    'UCODE': '21717',
                    'PCODE': 'NEWPASSWORD1',
                    'HUKEY': '0000000000' + s[0],
                    'ITEMUNPACK': {'PACK_QTY': s[1]},
                })
                print(f"Unpack ID Scrap: {outputScrap}")
                logger.info(f"Unpack ID Scrap: {outputScrap}")
            with open(work_order_name_f, newline='') as csvfile:
                wonamereader = csv.reader(
                    csvfile, delimiter=',', quotechar='|')
                for row in wonamereader:
                    work_order = '' .join(row)
            orderSAP = work_order

        dir = os.path.join(
            dataArchive+orderSAP)
        if not os.path.exists(dir):
            os.mkdir(dir)
        file_exist = os.path.exists(file_unpack_id)
        if file_exist:
            src_path = file_unpack_id
            dst_path = dataArchive + \
                orderSAP+"/unpack_id_"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)
        file_exist_scrap = os.path.exists(file_unpack_id_scrap)
        if file_exist_scrap:
            src_path = file_exist_scrap
            dst_path = dataArchive + \
                orderSAP+"/unpack_id_scrap_"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        connection.close()

    except CommunicationError:
        print("Could not connect to server.")
        logger.error("Could not connect to server.")
        logger.exception("Error!")
        raise
    except LogonError:
        print("Could not log in. Wrong credentials?")
        logger.error("Could not log in. Wrong credentials?")
        logger.exception("Error!")
        raise
    except (ABAPApplicationError, ABAPRuntimeError, ExternalRuntimeError):
        print("An error occurred.")
        logger.error("An error occurred.")
        logger.exception("Error!")
        raise


if __name__ == "__main__":
    main()
