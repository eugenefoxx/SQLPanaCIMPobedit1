from cmath import log
import configparser
import os

import pyrfc
import csv
from os.path import exists
from collections import defaultdict
from pyrfc import Connection
from datetime import datetime
import shutil

from pyrfc import ABAPApplicationError, ABAPRuntimeError, LogonError, CommunicationError, ExternalRuntimeError
import logging


def main():
    logger = logging.getLogger("order_info")
    logger.setLevel(logging.INFO)

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

    global wbs_el

    try:
        connection = pyrfc.Connection(**params_connection)
      #  result = connection.call('STFC_CONNECTION', REQUTEXT=u'Hello SAP!')
      #  print(result)
        logger.info("Connection to SAP RFC creating.")

        ttime = datetime.now()

        work_order = None
        orderSAP = None

        # work_order_name_f = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test1_work_order_name.csv"
        work_order_name_f = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test2_work_order_name.csv"

        with open(work_order_name_f, newline='') as csvfile:
            wonamereader = csv.reader(csvfile, delimiter=',', quotechar='|')
            for row in wonamereader:
                work_order = '' .join(row)
        if len(work_order) == 8:
            orderSAP = '0000' + work_order
        if len(work_order) == 7:
            orderSAP = '00000' + work_order
        orderSAPFolder = work_order
        # orderSAP = '00000' + work_order  # 000001000825 000001000836

        dataArchive = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_archive/"

        # wo_component = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test1_wo_component.csv"
        wo_component = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_spp_5/test2_wo_component.csv"

        # scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/scrap.csv"
        # scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test3_wo_component_scrap.csv"
        # scrap = "/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test_v2/test4_wo_component_scrap.csv"

        #infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data/info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/1000862_info_order.csv'
        infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test1_info_order.csv'
        # infoOrder = '/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/internal/pysaprfc/data_test/test2_info_order.csv'

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': str(orderSAP),  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']
        logger.info(f"Z_IEXT_PRODORD_INFO/RESITEMS: {order_info['RESITEMS']}")
        productsap = order_info['PRODUCT']
        logger.info(f"Z_IEXT_PRODORD_INFO/PRODUCT: {order_info['PRODUCT']}")
        print("productsap", productsap)
        print(order_info['RESITEMS'])
        print("_________inforesSAPorder___________")

        resSAPorder = [{'MATNR': sub['MATNR'], 'RSPOS': sub['RSPOS'], 'MATKL': sub['MATKL'], 'ERFMG': sub['ERFMG'], 'POSID': sub['POSID']}
                       for sub in ordp]

        #print("type", type(resSAPorder))
        print("resSAPorder:", resSAPorder)
        print("_________info___________")
        print(order_info)

        for i in resSAPorder:
            wbs_el = i['POSID']

        print("reserv number:", get_reserv_num(order_info))
        reserv = get_reserv_num(order_info)
        # new order - 1000836
        # записываем в массив список компонентов в заказе
        resmtlrs = [sub['MATNR'] for sub in ordp]
        print("resmtlrs: " + str(resmtlrs))

        # убираем, если есть полуфабрикат из листа
        arrPRODORD_INFO_Component = []
        for pf in resmtlrs:
            if not pf.__contains__('0000000000031'):
                arrPRODORD_INFO_Component.append(pf)

        # если заказ стадии 2, ставим полуфабрикату 1 стадии склад
        for rowMaterial in resSAPorder:
            if rowMaterial['MATNR'].__contains__('0000000000031'):
                chgMatr = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                    'UCODE': '21717',
                    'PCODE': 'NEWPASSWORD1',
                    'RESITEMS': [
                        {
                            u'LINE_ID': '2',
                            u'MATERIAL': rowMaterial['MATNR'],
                            u'PLANT': 'SL00',
                            u'STGE_LOC': '7813',
                            # u'BATCH': row['Lot'],
                            u'MOVE_TYPE': '261',
                            # u'ENTRY_QNT': row['SUM'],
                            u'ENTRY_UOM': 'ST',
                            u'ORDERID': orderSAP,
                            u'RESERV_NO': reserv,
                            u'RES_ITEM': rowMaterial['RSPOS'],
                            # u'WBS_ELEM': rowMaterial['POSID'],
                        }
                    ]
                })
                print(f"chgMatr {chgMatr}")
                logger.info(
                    f"Выставляем склад полуфабрикату в стадии 2: {chgMatr}")

        rowsPanaData = []
        with open(wo_component, newline='') as file:
            # csvreader = csv.reader(file, delimiter=',')
            csvreader = csv.DictReader(file, delimiter=',')
            # header = next(csvreader)
            for row in csvreader:
                rowsPanaData.append(row)
        print("rowsPanaData - ", rowsPanaData)
        for row in rowsPanaData:
            for c in resSAPorder:  # for c in resmtlrs:
                #    if c.__contains__('0000000000031'):
                #        print("31*", c)
                # проверка компонентов по Panacim на наличие в заказе
                if str('00000000000' + row['PART_NO']) == str(c['MATNR']):
                    #    if not c.__contains__('0000000000031'):
                    #    print('00000000000'+ row['PART_NO'])
                    print("This sap have in order", c['MATNR'], c['RSPOS'])

                    chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                        'UCODE': '21717',
                        'PCODE': 'NEWPASSWORD1',
                        'RESITEMS': [
                            {
                                u'LINE_ID': '2',
                                u'MATERIAL': c['MATNR'],
                                u'PLANT': 'SL00',
                                u'STGE_LOC': '7813',
                                u'BATCH': row['Lot'],
                                u'MOVE_TYPE': '261',
                                u'ENTRY_QNT': row['SUM'],
                                u'ENTRY_UOM': 'ST',
                                u'ORDERID': orderSAP,
                                u'RESERV_NO': reserv,
                                u'RES_ITEM': c['RSPOS'],
                                # u'WBS_ELEM': c['POSID'],
                            }
                        ]
                    })
                    print(chg)
                    logger.info(f"Z_IEXT_PRODORD_CHGRES 1: {chg}")
                    # logger.warning(parse_response(chg))
        #    elif str('00000000000' + row['PART_NO']) != str(c):
        #        print("not", c)

        # запись компонентов из Панасим в массив
        arrComponentFromPanaCIM = []
        with open(wo_component, newline='') as file:
            csvreader = csv.DictReader(file, delimiter=',')
            for row in csvreader:
                arrComponentFromPanaCIM.append('00000000000' + row['PART_NO'])
        print("arrComponentFromPanaCIM -", arrComponentFromPanaCIM)

        # повторно читаем информацию по заказу для добавления на дробление сап и партии
        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']
        logger.info(f"Z_IEXT_PRODORD_INFO/RESITEMS: {ordp}")

        # resSAPorder2 = [{'MATNR': sub['MATNR'], 'BDMNG': sub['BDMNG'], 'CHARG': sub['CHARG']} for sub in ordp]
        # resSAPorder2 = [{'PART_NO': sub['MATNR'], 'SUM': sub['BDMNG'], 'Lot': sub['CHARG']} for sub in ordp]

        resSAPorder2 = {}
        for i in ordp:
            resSAPorder2[i['MATNR']] = {
                "sum": i["BDMNG"], "lot": i["CHARG"], "wbs_elem": i["POSID"]}
        # resSAPorder2 = [{'Lot': sub['CHARG']} for sub in ordp]
        print("resSAPorder2:", resSAPorder2)
        #resarrComponentFromPanaCIM = []
        #print("rowsPanaData: ", rowsPanaData)

        # добавлем в SAP заказ компонент, согласно PanaCIM, если его там нет
        for i in rowsPanaData:  # arrComponentFromPanaCIM
            # if '00000000000' + i['PART_NO'] and i['Lot'] not in resSAPorder2:
            # if ('00000000000' + i['PART_NO'] and i['SUM'] and i['Lot']) not in resSAPorder2:
            print(str('00000000000' + i["PART_NO"]))
            if str('00000000000' + i["PART_NO"]) in resSAPorder2.keys():
                #  добавляем позицию исходя из списка Панасим
                if i["Lot"] != resSAPorder2['00000000000' + i["PART_NO"]]["lot"]:
                    print(f"Lot {i['Lot']} {i['PART_NO']} component will add")
                    addcomp = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                        'UCODE': '21717',
                        'PCODE': 'NEWPASSWORD1',
                        'RESITEMS': [
                            {
                                u'LINE_ID': '1',
                                u'MATERIAL': '00000000000' + i['PART_NO'],
                                u'PLANT': 'SL00',
                                u'STGE_LOC': '7813',
                                u'BATCH': i['Lot'],
                                u'MOVE_TYPE': '261',
                                u'ENTRY_QNT': i['SUM'],
                                u'ENTRY_UOM': 'ST',
                                u'ORDERID': orderSAP,
                                # u'WBS_ELEM': wbs_el,
                                #    u'RESERV_NO': reserv,
                                #    u'RES_ITEM': c['RSPOS'],
                            }
                        ]
                    })
                    print(addcomp)
                    logger.info(f"Z_IEXT_PRODORD_CHGRES 2: {addcomp}")
            else:
                # for c in resmtlrs:
                #  if '00000000000' + i['PART_NO'] != c['MATNR']:
                # print("component not have in sap_order:", i['PART_NO'], i['SUM'], i['Lot'])

                print(f"{i['PART_NO']} {i['Lot']} not found")
                addcompNotFound = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                    'UCODE': '21717',
                    'PCODE': 'NEWPASSWORD1',
                    'RESITEMS': [
                        {
                            'LINE_ID': '1',
                            'MATERIAL': '00000000000' + i['PART_NO'],
                            'PLANT': 'SL00',
                            'STGE_LOC': '7813',
                            'BATCH': i['Lot'],
                            'MOVE_TYPE': '261',
                            'ENTRY_QNT': i['SUM'],
                            'ENTRY_UOM': 'ST',
                            'ORDERID': orderSAP,
                            # 'WBS_ELEM': wbs_el,
                        }
                    ]
                })
                print(f"addcompNotFound: {addcompNotFound}")
                logger.info(f"addcompNotFound: {addcompNotFound}")
                # resarrComponentFromPanaCIM.append(i)

        # поиск компонентов, которые отсутствуют в cписке Panacim

        for i in resSAPorder:  # arrPRODORD_INFO_Component
            # print("remove:", str(i['MATNR']).removeprefix('00000000000'))
            # ch = str(i['MATNR']).removeprefix('00000000000')
            # print("ch:", ch)
            # arrComponentFromPanaCIM
            # 10502    Паяльные материалы    Материалы и сырье\Пайка\Паяльные материалы
            # 11204    Печатные платы    Материалы и сырье\Соединит.и изолир.комп.\Печатные платы
            if i['MATNR'] not in arrComponentFromPanaCIM and not i['MATNR'].__contains__(
                    '0000000000031') and i['MATKL'] != '11004' and i['MATKL'] != '10502' and i['ERFMG'] < 0:
                print("not in", i['MATNR'])
                zerocomp = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                    'UCODE': '21717',
                    'PCODE': 'NEWPASSWORD1',
                    'RESITEMS': [
                        {
                            u'LINE_ID': '6',
                            # u'MATERIAL': i['MATNR'],
                            u'PLANT': 'SL00',
                            # u'STGE_LOC': '7813',
                            # u'MOVE_TYPE': '261',
                            # u'ENTRY_QNT': '0.0',
                            # u'ENTRY_UOM': 'ST',
                            u'ORDERID': orderSAP,
                            u'RESERV_NO': reserv,
                            u'RES_ITEM': i['RSPOS'],
                            # u'WBS_ELEM': rowMaterial['POSID'],
                        }
                    ]
                })
                print(zerocomp)
                logger.info(f"Z_IEXT_PRODORD_CHGRES 3: {zerocomp}")

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        #  Вставка партии в добавленные строки
        sap_order = order_info['RESITEMS']
        print("Состояние проверки измененного заказа", sap_order)

        good_order = []
        bad_order = []
        rspos_black_list = []
        for i in sap_order:
            if not i['CHARG']:
                print(
                    f"Part number {i['MATNR']}, rspos {i['RSPOS']} does not have a lot")
                bad_order.append(i)
            else:
                good_order.append(i)
        print("bad_order", bad_order)
        print("rowsPanaData", rowsPanaData)
        for i in rowsPanaData:
            match = False
            for item in good_order:
                if item["MATNR"] == '00000000000' + i["PART_NO"] and item["CHARG"] == i["Lot"]:
                    match = True
                    break
                elif item["MATNR"] == '00000000000' + i["PART_NO"] and not item["CHARG"] == i["Lot"]:
                    match = False

            if not match:
                for item in bad_order:
                    if '00000000000' + i["PART_NO"] == item["MATNR"] and item['RSPOS'] not in rspos_black_list:
                        print({"Item of Bad_Order PART_NO": i['PART_NO'], "Lot": i['Lot'], "sum": i['SUM'],
                               "RSPOS": item['RSPOS']})
                        rspos_black_list.append(item['RSPOS'])

                        chg = connection.call('Z_IEXT_PRODORD_CHGRES', **{
                            'UCODE': '21717',
                            'PCODE': 'NEWPASSWORD1',
                            'RESITEMS': [
                                {
                                    u'LINE_ID': '2',
                                    u'MATERIAL': '00000000000' + i['PART_NO'],
                                    u'PLANT': 'SL00',
                                    u'STGE_LOC': '7813',
                                    u'BATCH': i['Lot'],
                                    u'MOVE_TYPE': '261',
                                    u'ENTRY_QNT': i['SUM'],
                                    u'ENTRY_UOM': 'ST',
                                    u'ORDERID': orderSAP,
                                    u'RESERV_NO': reserv,
                                    u'RES_ITEM': item['RSPOS'],
                                    # u'WBS_ELEM': item['POSID'],
                                }
                            ]
                        })
                        print(chg)
                        logger.info(
                            f"Z_IEXT_PRODORD_CHGRES: Вставка партии в добавленные строки: {chg}")
                        break

                    elif item['RSPOS'] in rspos_black_list:
                        pass
                    else:
                        print(
                            f"Part number {i['PART_NO']} was not found in SAPOrder")

        order_info = connection.call('Z_IEXT_PRODORD_INFO', **{
            'AUFNR': orderSAP,  # 000001000825
            'UCODE': '21717',
            'PCODE': 'NEWPASSWORD1',
        }
        )
        ordp = order_info['RESITEMS']
        print("Состояние проверки измененного заказа 2", ordp)

        dir = os.path.join(dataArchive+orderSAPFolder)
        if not os.path.exists(dir):
            os.mkdir(dir)

        # чтение файла по скрапу
        # rowsScrap = []
        # file_existScrap = os.path.exists(scrap)
        # if file_existScrap:
        #    with open(scrap, newline='') as scrapFile:
        #        csvreader = csv.DictReader(scrapFile, delimiter=',')
        #        for row in csvreader:
        #            rowsScrap.append(row)

            # добавляем в SAP заказ компоненты из списка по scrap
        #    for i in rowsScrap:
        #        addScrap = connection.call('Z_IEXT_PRODORD_CHGRES', **{
        #            'UCODE': '21717',
        #            'PCODE': 'NEWPASSWORD1',
        #            'RESITEMS': [
        #                {
        #                    'LINE_ID': '1',
        #                    'MATERIAL': '00000000000' + i['PART_NO'],
        #                    'PLANT': 'SL00',
        #                    'STGE_LOC': '7813',
        #                    'BATCH': i['Lot'],
        #                    'MOVE_TYPE': '261',
        #                    'ENTRY_QNT': i['SUM'],
        #                    'ENTRY_UOM': 'ST',
        #                    'ORDERID': orderSAP,
        #                }
        #            ]
        #        })
        #        print(f"add srap: {addScrap}")
        #        logger.info(f"add srap: {addScrap}")

        #    src_path = scrap
        #    dst_path = dataArchive + orderSAP+"/scrap"+str(ttime)+".csv"
        #    shutil.move(src_path, dst_path)

        file_existwocomp = os.path.exists(wo_component)
        if file_existwocomp:
            src_path = wo_component
            dst_path = dataArchive + orderSAPFolder + \
                "/wo_component_"+str(ttime)+".csv"
            shutil.move(src_path, dst_path)

        file_exist_wo_name = os.path.exists(work_order_name_f)
        if file_exist_wo_name:
            src_path = work_order_name_f
            dst_path = dataArchive + orderSAPFolder + \
                "/work_order_name_"+str(ttime)+".csv"
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
    except (ABAPApplicationError, ABAPRuntimeError):
        print("An error occurred.")
        logger.error("An error occurred.")
        logger.exception("Error!")
        raise


def get_reserv_num(dict_value):
    if not dict_value:
        return "Ответ не получен"
    else:
        for key, value in dict_value.items():
            if key == 'RESITEMS':
                res = value[0].get('RSNUM', '')
                return res


def parse_response(dict_value):
    if not dict_value:
        return "Ответ не получен"
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
            #  print("NOK")
            # elif value[0].get('NUMBER', '') == '100':
            # print("OK")
            # else:
            #    return "Ответ не получен"


# def parse_response(dict_value):
#    for key, value in dict_value.items():
#        if key == 'RETURN':
#            if value[0].get('TYPE', '') == 'E':
#                return 'Error: ' + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'I':
#                return "Infomation: " + str(value[0].get('MESSAGE', ''))
#            elif value[0].get('TYPE', '') == 'W':
#                return "Warning: " + str(value[0].get('MESSAGE', ''))
            # elif value[0].get('NUMBER', '') == '469':
            # print("NOK")
            # elif value[0].get('NUMBER', '') == '100':
            # print("OK")
#            else:
#                return "Ответ не получен"


if __name__ == "__main__":
    main()
