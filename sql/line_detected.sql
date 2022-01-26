-- определяем MIX_NAME изделия с группированием по времени модификации и определяем по MIX_NAME [ROUTE_ID]
SELECT TOP 1 [PRODUCT_ID]
      ,[ROUTE_ID]
      ,[MIX_NAME]
      ,[SETUP_ID]
      ,[LDF_FILE_NAME]
      ,[MACHINE_FILE_NAME]
      ,[SETUP_VALID_FLAG]
      ,[LAST_MODIFIED_TIME]
      ,[DOS_FILE_NAME]
      ,[MODEL_STRING]
      ,[TOP_BOTTOM]
      ,[PT_GROUP_NAME]
      ,[PT_LOT_NAME]
      ,[PT_MC_FILE_NAME]
      ,[PT_DOWNLOADED_FLAG]
      ,[PT_NEEDS_DOWNLOAD]
      ,[SUB_PARTS_FLAG]
      ,[BARCODE_SIDE]
      ,[CYCLE_TIME]
      ,[IMPORT_SOURCE]
      ,[MODIFIED_IMPORT_SOURCE]
      ,[THEORETICAL_XOVER_TIME]
      ,[PUBLISH_MODE]
      ,[PCB_NAME]
      ,[MASTER_MJS_ID]
      ,[LED_VALID_FLAG]
      ,[DGS_PPD_VALID_FLAG]
      ,[ACTIVE_MODEL_STRING]
      ,[REGISTERED_PCB_NAME]
FROM [PanaCIM].[dbo].[product_setup]
WHERE PRODUCT_ID = '3082'
order by LAST_MODIFIED_TIME desc

-- определяем по [ROUTE_ID] - [ROUTE_NAME], это определение линии
SELECT TOP 1000 [ROUTE_ID]
      ,[ROUTE_NAME]
      ,[HOST_NAME]
      ,[DOS_LINE_NO]
      ,[FLOW_DIRECTION]
      ,[VALID_FLAG]
      ,[SUBIMPORT_PATH]
      ,[STAND_ALONE]
      ,[ROUTE_STARTUP]
      ,[LNB_HOST_NAME]
      ,[ROUTE_ABBR]
      ,[DGS_LINE_ID]
      ,[DGS_IMPORT_MODE]
      ,[MGMT_UPLOAD_TYPE]
      ,[SUB_PART_IMPORT_SRC]
      ,[NAVI_IMPORT_MODE]
      ,[RESTRICTED_COMPONENTS_ENABLED]
      ,[SEPARATE_NETWORK_IP]
      ,[SEPARATE_NETWORK_ENABLED]
      ,[PUBLISH_MODE]
      ,[LINKED_TO_PUBLISH]
      ,[PUBLISH_ROUTE_ID]
      ,[DISABLE_TRAY_PART_SCAN]
      ,[ENABLE_TRAY_INTERLOCK]
      ,[ALLOW_DELETE]
      ,[BMX_ZONE_ID]
      ,[BMX_STORAGE_UNIT_ID]
      ,[BMX_DEDICATION_TYPE]
      ,[PT200_LINE_ID]
      ,[DGS_SERVER_ID]
      ,[ERP_ROUTE_ID]
      ,[ENABLE_TRAY_RFID]
      ,[TRAY_MODE_MOD1]
      ,[TRAY_MODE_MOD2]
      ,[LMV2_ID]
      ,[ILNB_HOST_NAME]
      ,[ARDCE_LINE_NO]
  FROM [PanaCIM].[dbo].[routes]