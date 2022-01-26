/****** Script for SelectTopNRows command from SSMS  ******/
-- посмотреть пару job_id и setup_id
SELECT TOP 5000 [JOB_ID]
      ,[SETUP_ID]
FROM [PanaCIM].[dbo].[job_products]
order by [SETUP_ID] desc

-- посмотреть список ордеров по изделиям
SELECT TOP 100000 [WORK_ORDER_ID]
      ,[WORK_ORDER_NAME]
      ,[LOT_SIZE]
      ,[JOB_ID]
      ,[MASTER_WORK_ORDER_ID]
      ,[COMMENTS]
FROM [PanaCIM].[dbo].[work_orders]
order by [JOB_ID] desc

-- get PCB_NAME
SELECT [PCB_NAME]
FROM [PanaCIM].[dbo].[product_setup]
WHERE PRODUCT_ID = '3410'
group by [PCB_NAME]

-- PT_LOT_NAME в [PanaCIM].[dbo].[product_setup]
-- это стороны PCB_NAME

-- get NPM name for recipte
SELECT TOP 1000 [PRODUCT_ID]
      ,[PRODUCT_NAME]
      ,[DOS_PRODUCT_NAME]
      ,[PATTERNS_PER_PANEL]
      ,[PANEL_WIDTH]
      ,[PANEL_LENGTH]
      ,[PANEL_THICKNESS]
      ,[CAMERA_XAXIS_TOP]
      ,[CAMERA_YAXIS_TOP]
      ,[CAMERA_XAXIS_BOTTOM]
      ,[CAMERA_YAXIS_BOTTOM]
      ,[TOOLING_PIN_DISTANCE]
      ,[BARCODES_PER_PANEL]
      ,[PRODUCT_VALID_FLAG]
      ,[TOOLING_PIN]
      ,[CONVEYOR_SPEED]
      ,[USE_BRD_FILE]
      ,[BASE_PRODUCT_ID]
      ,[PATTERN_COMBINATIONS_PER_PANEL]
      ,[PATTERN_TYPES_PER_PANEL]
  FROM [PanaCIM].[dbo].[product_data]
  where [PRODUCT_ID] = '2501'