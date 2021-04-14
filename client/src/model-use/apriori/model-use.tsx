import { useEffect } from 'react';
import { changePageTitle, changeSelectedMenuKey } from '../../app.slice';
import { FetchState } from '../../domain/fetch-state';
import { useAppDispatch, useAppSelector } from '../../hooks';
import { changeLoadingExamples, fetchExampleRules, fetchRecommendations } from './model-use.slice';
import { Divider, Space, Spin, Typography } from 'antd';
import { ArrowRightOutlined } from '@ant-design/icons';
import './model-use.scss'
import Search from 'antd/lib/input/Search';
import { MenuKey } from '../../domain/menu-key.enum';

const { Title } = Typography;

export default function ModelUse() {
  const PAGE_TITLE = 'Model Usage';
  const dispatch = useAppDispatch();
  const loadingExamples = useAppSelector(state => state.modelUse.loadingExamples);
  const exampleRules = useAppSelector(state => state.modelUse.exampleRules);
  const loadingRecommendations = useAppSelector(state => state.modelUse.loadingRecommendations);
  const pageTitle = useAppSelector(s => s.app.pageTitle);
  const recommendations = useAppSelector(state => state.modelUse.recommendations);

  useEffect(() => {
    if (pageTitle !== PAGE_TITLE) {
      // this triggers a full re-render
      dispatch(changePageTitle('Model Usage'));
      dispatch(changeSelectedMenuKey(MenuKey.ModelUsage));
      dispatch(changeLoadingExamples(FetchState.idle));
    } else if (loadingExamples === FetchState.idle) {
      dispatch(fetchExampleRules());
    }
  })

  function handleTest(v: string) {
    dispatch(fetchRecommendations(v.split(',')));
  }

  const noDataContainer = (
    <Title level={4} type="secondary" >No Rules found.</Title>
  )

  const recommendationTestContainer = (
    <>
      <Title level={4} style={{ marginBottom: '10px' }}>Test Recommendations</Title>
      <Search
        placeholder="1234,5434,4563"
        allowClear
        enterButton="Send"
        size="large"
        onSearch={handleTest}
        style={{ width: '50%' }}
      />
      {
        loadingRecommendations === FetchState.loading &&
        <div className="loader-positioner">
          <Space size="middle">
            <Spin size="large" />
          </Space>
        </div>
      }
      {
        loadingRecommendations === FetchState.succeeded &&
        <div style={{ marginTop: '20px', marginLeft: '20px' }}>
          <Title level={5} style={{ marginBottom: '5px' }}>Recommended Items:</Title>
          {
            recommendations.map(r =>
              <div key={r.prediction} className="recommendation">
                <div>Item: {r.prediction}</div>
                <div>Frequency: {r.frequency}</div>
                <div>Confidence: {Math.round(r.confidence * 100)}%</div>
              </div>
            )
          }
        </div>
      }
    </>
  )

  const exampleRulesContainer = (
    <>
      <Title level={4} style={{ marginBottom: '10px' }}>Example Rules</Title>
      {
        loadingExamples === FetchState.loading &&
        <div className="loader-positioner">
          <Space size="middle">
            <Spin size="large" />
          </Space>
        </div>
      }
      {
        loadingExamples === FetchState.succeeded &&
        exampleRules.map(r =>
          <div key={r.baseItemSet + r.prediction} className="rule">
            <div>{r.baseItemSet}</div>
            <div><ArrowRightOutlined /></div>
            <div>{r.prediction}</div>
          </div>
        )
      }
    </>
  )

  return (
    <>
      {
        loadingExamples === FetchState.succeeded && exampleRules?.length === 0 &&
        noDataContainer
      }
      {
        loadingExamples === FetchState.succeeded && exampleRules?.length > 0 &&
        recommendationTestContainer
      }
      {
        (loadingExamples !== FetchState.succeeded || exampleRules?.length > 0) &&
        <>
          <Divider />
          {exampleRulesContainer}
        </>
      }
    </>
  )
}