/* поиск последнего закрывшегося job_id */
/****** Script for SelectTopNRows command from SSMS  ******/
IF OBJECT_ID('dbo.GetLastJobId') IS NOT NULL DROP FUNCTION dbo.GetLastJobId;
GO

CREATE FUNCTION dbo.GetLastJobId
RETURNS INT

SELECT TOP 1 [JOB_ID]
      FROM [PanaCIM].[dbo].[job_history] where CLOSING_TYPE = '0' order by END_TIME desc;
      GO
