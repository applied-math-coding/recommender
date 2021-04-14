import { Space, Spin } from 'antd';
import { useEffect } from 'react';
import { changePageTitle, changeSelectedMenuKey } from '../../app.slice';
import { FetchState } from '../../domain/fetch-state';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { changeStatsLoading, fetchModelStats } from './model-stats.slice';
import "./model-stats.scss";
import { Histogram } from '@ant-design/charts';
import { listToChartData } from '../../common/utils';
import { Typography } from 'antd';
import { MenuKey } from '../../domain/menu-key.enum';

const { Title } = Typography;

export default function ModelStats() {
  const PAGE_TITLE = 'Model Statistics';
  const statsLoading = useAppSelector(s => s.modelStats.statsLoading);
  const ruleStatistics = useAppSelector(s => s.modelStats.ruleStatistics);
  const pageTitle = useAppSelector(s => s.app.pageTitle);
  const dispatch = useAppDispatch();

  useEffect(() => {
    if (pageTitle !== PAGE_TITLE) {
      // this triggers a full re-render
      dispatch(changePageTitle(PAGE_TITLE));
      dispatch(changeSelectedMenuKey(MenuKey.ModelStatistics));
      dispatch(changeStatsLoading(FetchState.idle));
    } else if (statsLoading === FetchState.idle) {
      dispatch(fetchModelStats());
    }
  });

  const noDataContainer = (
    <Title level={4} type="secondary">No Model found.</Title>
  )

  const loader = (
    <div className="loader-positioner">
      <Space size="middle">
        <Spin size="large" />
      </Space>
    </div>
  )

  const statsContainer = (
    <>
      <Title level={4} style={{ textAlign: 'center' }}>Distributions of Item-Frequencies</Title>
      {
        statsLoading === FetchState.succeeded &&
        ruleStatistics.map(rs =>
          <div key={rs.baseLength} className="chart-container">
            <Title level={4} style={{ marginBottom: '10px' }}>Rule-Length: {rs.baseLength}</Title>
            <Histogram binWidth={1} binField={'value'} data={listToChartData(rs.frequencies)} />
          </div>
        )
      }
    </>
  )

  return (
    <>
      {
        statsLoading === FetchState.loading &&
        loader
      }
      {
        statsLoading === FetchState.succeeded && ruleStatistics?.length === 0 &&
        noDataContainer
      }
      {
        statsLoading === FetchState.succeeded && ruleStatistics?.length > 0 &&
        statsContainer
      }
    </>
  )
}
