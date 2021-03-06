import configparser
import csv
import os
from os.path import exists
from datetime import datetime
import shutil
from decimal import *

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
    global paramsATHDRLEVELS, paramsGOODSMOVEMENTS, SAP_ORDER, SAP_ORDER_Number
    # try:
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

    global wbs_el

    ttime = datetime.now()
    try:

        # while True:
        connection = pyrfc.Connection(**params_connection)
       # result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
       # print(result)
        logger.info("Connection to SAP RFC creating.")
       # resultTime = connection.call(
       #     'WEEK_GET_FIRST_DAY', **{'WEEK': '201825'})
       # print(resultTime)
        dataArchive = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_archive/"

        # infoOrder = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test1_info_order.csv"
        infoOrder = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test2_info_order.csv"

        # scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test1_wo_component_scrap.csv"
        scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test2_wo_component_scrap.csv"
        orderSAP = None
        rowsinfoOrder = []
        with open(infoOrder, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                rowsinfoOrder.append(row)
        for row in rowsinfoOrder:
            if len(row['WO']) == 8:
                orderSAP = '0000' + row['WO']
            if len(row['WO']) == 7:
                orderSAP = '00000' + row['WO']
            paramsATHDRLEVELS = [

                {
                    # '000001000825s',  # 000001000825
                    'ORDERID': orderSAP,  # '00000' + row['WO'],
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
            if len(i) == 8:
                SAP_ORDER = '0000' + i
            if len(i) == 7:
                SAP_ORDER = '00000' + i
            #SAP_ORDER = i
            SAP_ORDER_Number = i
        # print("sap order", rowsinfoOrder['WO'])
        print(SAP_ORDER)

        # infoMaterialOrder = []
        # with open(infomaterialOrder, newline='') as fileMaterial:
        #    csvreader = csv.DictReader(fileMaterial, delimiter=',')
        #    for row in csvreader:
        #        infoMaterialOrder.append(row)

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': SAP_ORDER,  # '00000' + SAP_ORDER,
            # '000001000836', # str('00000' + str(sapORDER)),  # '00000' + row['WO'],  # orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        productsap = order_info['PRODUCT']
        print("productsap", productsap)
        logger.info(f"productsap: {productsap}")
        sap_order = order_info['RESITEMS']
        print("sap_order", sap_order)
        logger.info(f"sap_order: {sap_order}")

        wbs_l = [{'POSID': sub['POSID']} for sub in sap_order]
        for i in wbs_l:
            wbs_el = i['POSID']

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
                # '00000' + row['WO'],  # '000001000825',
                'ORDERID': SAP_ORDER,
                'REF_DOC_IT': '0001',
            })

        for matr in sap_order:
            if matr['MATNR'].__contains__('0000000000031') and wbs_l != '':
                paramsGOODSMOVEMENTS.append({
                    'MATERIAL': matr['MATNR'],  # '000000000002003411',
                    'ENTRY_QNT': matr['ERFMG'],  # '1',
                    'ENTRY_UOM': matr['MEINS'],  # 'ST',
                    'STGE_LOC': '7813',  # matr['LGORT']
                    'BATCH': matr['CHARG'],  # '1000001747',
                    'MOVE_TYPE': '261',
                    # Индикатор особого запаса matr['SOBKZ']
                    'SPEC_STOCK': 'Q',
                    'WBS_ELEM': matr['POSID'],  # СПП-элемент
                    'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                    # '0000031904',  # Номер резерва
                    'RESERV_NO': matr['RSNUM'],
                    # '0020',  # Номер позиции резерва
                    'RES_ITEM': matr['RSPOS'],
                    'PLANT': 'SL00',  # matr['WERKS']
                    'ORDERID': SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                    'WITHDRAWN': 'X',  # фиксированное значение
                    'REF_DOC_IT': '0001',  # фиксированное значение
                })
            else:
                paramsGOODSMOVEMENTS.append({
                    'MATERIAL': matr['MATNR'],  # '000000000002003411',
                    'ENTRY_QNT': matr['ERFMG'],  # '1',
                    'ENTRY_UOM': matr['MEINS'],  # 'ST',
                    'STGE_LOC': '7813',  # matr['LGORT']
                    'BATCH': matr['CHARG'],  # '1000001747',
                    'MOVE_TYPE': '261',
                    'SPEC_STOCK': matr['SOBKZ'],  # Индикатор особого запаса
                    'WBS_ELEM': matr['POSID'],  # СПП-элемент
                    'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                    # '0000031904',  # Номер резерва
                    'RESERV_NO': matr['RSNUM'],
                    # '0020',  # Номер позиции резерва
                    'RES_ITEM': matr['RSPOS'],
                    'PLANT': 'SL00',  # matr['WERKS']
                    'ORDERID': SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                    'WITHDRAWN': 'X',  # фиксированное значение
                    'REF_DOC_IT': '0001',  # фиксированное значение
                })

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
                if wbs_l != '':
                    paramsGOODSMOVEMENTS.append({
                        # '000000000002003411',
                        'MATERIAL': '00000000000'+i['PART_NO'],
                        'ENTRY_QNT': i['SUM'],  # '1',
                        'ENTRY_UOM': 'ST',
                        'STGE_LOC': '7813',
                        'BATCH': i['Lot'],  # '1000001747',
                        'MOVE_TYPE': 'Z61',
                        'SPEC_STOCK': 'Q',  # Индикатор особого запаса
                        'WBS_ELEM': wbs_el,  # СПП-элемент
                        'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                        # '0000031904',  # Номер резерва
                        # 'RESERV_NO': matr['RSNUM'],
                        # '0020',  # Номер позиции резерва
                        # 'RES_ITEM': matr['RSPOS'],
                        'PLANT': 'SL00',
                        'ORDERID': SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                        'WITHDRAWN': 'X',  # фиксированное значение
                        'REF_DOC_IT': '0001',  # фиксированное значение
                    })
                else:
                    paramsGOODSMOVEMENTS.append({
                        # '000000000002003411',
                        'MATERIAL': '00000000000'+i['PART_NO'],
                        'ENTRY_QNT': i['SUM'],  # '1',
                        'ENTRY_UOM': 'ST',
                        'STGE_LOC': '7813',
                        'BATCH': i['Lot'],  # '1000001747',
                        'MOVE_TYPE': 'Z61',
                        'SPEC_STOCK': '',  # Индикатор особого запаса
                        'WBS_ELEM': wbs_el,  # СПП-элемент
                        'NO_MORE_GR': '',  # = 'X' если конечное подтверждение
                        # '0000031904',  # Номер резерва
                        # 'RESERV_NO': matr['RSNUM'],
                        # '0020',  # Номер позиции резерва
                        # 'RES_ITEM': matr['RSPOS'],
                        'PLANT': 'SL00',
                        'ORDERID': SAP_ORDER,  # '000001000825', '00000' + SAP_ORDER
                        'WITHDRAWN': 'X',  # фиксированное значение
                        'REF_DOC_IT': '0001',  # фиксированное значение
                    })

                print(
                    f"add srap in paramsGOODSMOVEMENTS: {paramsGOODSMOVEMENTS}")

            src_path = scrap
            dst_path = dataArchive + SAP_ORDER_Number + \
                "/scrap_"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        print(f"add in paramsGOODSMOVEMENTS: {paramsGOODSMOVEMENTS}")
        logger.info(f"add in paramsGOODSMOVEMENTS: {paramsGOODSMOVEMENTS}")
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

        file_exist = os.path.exists(infoOrder)
        if file_exist:
            src_path = infoOrder
            dst_path = dataArchive + SAP_ORDER_Number + \
                "/info_order_"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        print(outputtedorder)
        logger.info(f"outputtedorder: {outputtedorder}")

        parse_response(outputtedorder)
        # logger.warning(parse_response(
        #    f"outputtedorder parse: {outputtedorder}"))
        logger.warning(
            f"outputtedorder parse {parse_response(outputtedorder)}")
        print(parse_response(outputtedorder))

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


# def parse_response(dict_value):
#    if not dict_value:
#        return "Сообщения нет"
#    else:
#        for key, value in dict_value.items():
#            if key == 'RETURN':
#                if value.get('TYPE', '') == 'E':
#                    return 'Error: ' + str(value.get('MESSAGE', ''))
#                elif value.get('TYPE', '') == 'I':
#                    return "Infomation: " + str(value.get('MESSAGE', ''))
#                elif value.get('TYPE', '') == 'W':
#                    return "Warning: " + str(value.get('MESSAGE', ''))


def parse_response(dict_value):
    dv = dict_value['RETURN']
    if dv == []:
        return "RETURN: Сообщения нет"
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
