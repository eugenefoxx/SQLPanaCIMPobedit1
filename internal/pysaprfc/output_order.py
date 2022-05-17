import configparser
import csv
import os
from os.path import exists
from datetime import datetime
import shutil

import pyrfc

from pyrfc import Connection

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError
import logging


# выпуск изделия
def main():
    # log_file = 'logging.log'
    # fl = open(log_file, 'a+')
    # fl.close()
    # logging.basicConfig(filename="logging.log", level=logging.INFO)
    global paramsATHDRLEVELS, paramsGOODSMOVEMENTS, SAP_ORDER
    try:
        logger = logging.getLogger("output_order")
        logger.setLevel(logging.INFO)

        # create the logging file handler
        fh = logging.FileHandler(
            "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/pyrfc_logging.log")
        formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
        fh.setFormatter(formatter)

        # add handler to logger object
        logger.addHandler(fh)

        logger.info("Parsing cfg")
        config = configparser.ConfigParser()
        config.read(
            "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/sapnwrfc.cfg")
        config.sections()
        params_connection = config['connection']
        logger.info(f"Connecting to SAP RFC...")

        ttime = datetime.now()
    # try:

        # while True:
        connection = pyrfc.Connection(**params_connection)
       # result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
       # print(result)
        logger.info("Connection to SAP RFC creating.")
       # resultTime = connection.call(
       #     'WEEK_GET_FIRST_DAY', **{'WEEK': '201825'})
       # print(resultTime)
        dataArchive = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_archive/"
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/1000862_info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_info_order.csv'
        infoOrder = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_info_order.csv"

        scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_wo_component_scrap.csv"
        # scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_wo_component_scrap.csv"

        # infoOrder = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_info_order.csv"
        # 'info_material_order.csv'
        # infomaterialOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/wo_component.csv'
        # infomaterialOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/wo_component_1000862.csv'
        # infomaterialOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_wo_component.csv'
        # infomaterialOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_2_wo_component.csv'

        rowsinfoOrder = []
        with open(infoOrder, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                rowsinfoOrder.append(row)
        for row in rowsinfoOrder:
            paramsATHDRLEVELS = [

                {
                    # '000001000825s',  # 000001000825
                    'ORDERID': '00000' + row['WO'],
                    'YIELD': row['Qty'],  # '1.00',
                    'POSTG_DATE': row['Date'],  # '20220211',
                    'FIN_CONF': '',
                    'CLEAR_RES': '',
                }
            ]
        print("rowsinfoOrder:", type(rowsinfoOrder))
        sapORDER = [sub['WO'] for sub in rowsinfoOrder]
        print(str(sapORDER))
        for i in sapORDER:
            print("i", i)
            SAP_ORDER = i
        # print("sap order", rowsinfoOrder['WO'])
        print(SAP_ORDER)

        # infoMaterialOrder = []
        # with open(infomaterialOrder, newline='') as fileMaterial:
        #    csvreader = csv.DictReader(fileMaterial, delimiter=',')
        #    for row in csvreader:
        #        infoMaterialOrder.append(row)

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': '00000' + SAP_ORDER,  # '00000' + SAP_ORDER,
            # '000001000836', # str('00000' + str(sapORDER)),  # '00000' + row['WO'],  # orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        productsap = order_info['PRODUCT']
        print("productsap", productsap)
        logger.info(f"{productsap}")
        sap_order = order_info['RESITEMS']
        print("sap_order", sap_order)
        logger.info(f"{sap_order}")

        # breakpoint()
        paramsGOODSMOVEMENTS = []
        for row in rowsinfoOrder:
            paramsGOODSMOVEMENTS.append({
                'MATERIAL': productsap,  # '000000000003100302',
                'PLANT': 'SL00',
                'STGE_LOC': '7813',
                'MOVE_TYPE': '131',
                'ENTRY_QNT': row['Qty'],  # '1',
                'ENTRY_UOM': 'ST',
                'ORDERID': '00000' + row['WO'],  # '000001000825',
                'REF_DOC_IT': '0001',
            })

        print("Added first record in paramsGOODSMOVEMENTS", paramsGOODSMOVEMENTS)

        for matr in sap_order:
            paramsGOODSMOVEMENTS.append({
                'MATERIAL': matr['MATNR'],  # '000000000002003411',
                'ENTRY_QNT': matr['ERFMG'],  # '1',
                'ENTRY_UOM': 'ST',
                'STGE_LOC': '7813',
                'BATCH': matr['CHARG'],  # '1000001747',
                'MOVE_TYPE': '261',
                'SPEC_STOCK': '',  # Индикатор особого запаса
                'WBS_ELEM': '',  # СПП-элемент
                'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                'RESERV_NO': matr['RSNUM'],  # '0000031904',  # Номер резерва
                'RES_ITEM': matr['RSPOS'],  # '0020',  # Номер позиции резерва
                'PLANT': 'SL00',
                'ORDERID': '00000' + SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                'WITHDRAWN': 'X',  # фиксированное значение
                'REF_DOC_IT': '0001',  # фиксированное значение
            })
        print("Added second and after records in paramsGOODSMOVEMENTS",
              paramsGOODSMOVEMENTS)

        # добавление скрапа из файла
        rowsScrap = []
        file_existScrap = os.path.exists(scrap)
        if file_existScrap:
            with open(scrap, newline='') as scrapFile:
                csvreader = csv.DictReader(scrapFile, delimiter=',')
                for row in csvreader:
                    rowsScrap.append(row)

            # добавляем в SAP заказ компоненты из списка по scrap
            for i in rowsScrap:
                paramsGOODSMOVEMENTS.append({
                    # '000000000002003411',
                    'MATERIAL': '00000000000'+i['PART_NO'],
                    'ENTRY_QNT': i['SUM'],  # '1',
                    'ENTRY_UOM': 'ST',
                    'STGE_LOC': '7813',
                    'BATCH': i['Lot'],  # '1000001747',
                    'MOVE_TYPE': 'Z61',
                    'SPEC_STOCK': '',  # Индикатор особого запаса
                    'WBS_ELEM': '',  # СПП-элемент
                    'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                    # '0000031904',  # Номер резерва
                    # 'RESERV_NO': matr['RSNUM'],
                    # '0020',  # Номер позиции резерва
                    # 'RES_ITEM': matr['RSPOS'],
                    'PLANT': 'SL00',
                    'ORDERID': '00000' + SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                    'WITHDRAWN': 'X',  # фиксированное значение
                    'REF_DOC_IT': '0001',  # фиксированное значение
                })

                print(
                    f"add srap in paramsGOODSMOVEMENTS: {paramsGOODSMOVEMENTS}")
                logger.info(
                    f"add srap in paramsGOODSMOVEMENTS: {paramsGOODSMOVEMENTS}")

            src_path = scrap
            dst_path = dataArchive + SAP_ORDER+"/scrap"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        # breakpoint()
        # for row in rowsinfoOrder:
        #     for matr in sap_order:  # infoMaterialOrder:
        #         paramsGOODSMOVEMENTS = [
        #             {
        #                 'MATERIAL': productsap,  # '000000000003100302',
        #                 'PLANT': 'SL00',
        #                 'MOVE_TYPE': '131',
        #                 'ENTRY_QNT': row['Qty'],  # '1',
        #                 'ENTRY_UOM': 'ST',
        #                 'ORDERID': '00000' + row['WO'],  # '000001000825',
        #                 'REF_DOC_IT': '0001',
        #             },
        #             {
        #                 'MATERIAL': matr['MATNR'],  # '000000000002003411',
        #                 'ENTRY_QNT': matr['ERFMG'],  # '1',
        #                 'ENTRY_UOM': 'ST',
        #                 'STGE_LOC': '7813',
        #                 'BATCH': matr['CHARG'],  # '1000001747',
        #                 'MOVE_TYPE': '261',
        #                 'SPEC_STOCK': '',  # Индикатор особого запаса
        #                 'WBS_ELEM': '',  # СПП-элемент
        #                 'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
        #                 'RESERV_NO': matr['RSNUM'],  # '0000031904',  # Номер резерва
        #                 'RES_ITEM': matr['RSPOS'],  # '0020',  # Номер позиции резерва
        #                 'PLANT': 'SL00',
        #                 'ORDERID': '00000' + SAP_ORDER,  # '000001000825',
        #                 'WITHDRAWN': 'X',  # фиксированное значение
        #                 'REF_DOC_IT': '0001',  # фиксированное значение
        #             },
        #
        #         ]

        outputtedorder = connection.call('Z_IEXT_PRODORDCONF_CREATE_HDR', **{
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
            'ATHDRLEVELS': paramsATHDRLEVELS,
            'GOODSMOVEMENTS': paramsGOODSMOVEMENTS,
        }
        )

        print(outputtedorder)
        logger.info(f"{outputtedorder}")

        parse_response(outputtedorder)
        logger.warning(parse_response(f"{outputtedorder}"))
        print(parse_response(outputtedorder))

        file_exist = os.path.exists(infoOrder)
        if file_exist:
            src_path = infoOrder
            dst_path = dataArchive + SAP_ORDER+"/info_order"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        connection.close()
    # except KeyError:
    #    logger.error.

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
    except (ABAPApplicationError, ABAPRuntimeError):
        print("An error occurred.")
        logger.error("An error occurred.")
        logger.exception("Error!")
        raise


# def parse_response(dict_value):
#    for key, value in dict_value.items():
#        if key == 'RETURN':
#            if value[0].get('TYPE', '') == 'E':
#                return 'Error: ' + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'I':
#                return "Infomation: " + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'W':
#                return "Warning: " + str(value[0].get('MESSAGE', ''))
#            # elif value[0].get('NUMBER', '') == '469':
#            # print("NOK")
#            # elif value[0].get('NUMBER', '') == '100':
#            # print("OK")
#            else:
#                return "Ответ не получен"


def parse_response(dict_value):
    if not dict_value:
        return "Сообщения нет"
    else:
        for value in dict_value['RETURN']:
            # if key == 'RETURN':
            if value.get('TYPE', '') == 'E':
                return 'Error: ' + str(value.get('MESSAGE', ''))
            elif value.get('TYPE', '') == 'I':
                return "Infomation: " + str(value.get('MESSAGE', ''))
            elif value.get('TYPE', '') == 'W':
                return "Warning: " + str(value.get('MESSAGE', ''))

            # elif value[0].get('NUMBER', '') == '469':
            # print("NOK")
            # elif value[0].get('NUMBER', '') == '100':
            # print("OK")
            # else:
            #    return "Ответ не получен"


if __name__ == "__main__":
    main()
