import { Menu } from 'antd';
import Layout, { Content, Footer, Header } from 'antd/lib/layout/layout';
import React, { Suspense } from 'react';
import { Link, BrowserRouter as Router, Route, Switch, Redirect } from 'react-router-dom';
import './App.scss';
import { Loading } from './common/loading';
import { useAppDispatch, useAppSelector } from './hooks';
import { Typography } from 'antd';
import { MenuKey } from './domain/menu-key.enum';
import { changeSelectedMenuKey } from './app.slice';

function App() {
  const DataImport = React.lazy(() => import('./data-import/data-import'));
  const ModelStats = React.lazy(() => import('./model-stats/apriori/model-stats'));
  const CosineModelStats = React.lazy(() => import('./model-stats/cosine/cosine-model-stats'));
  const ModelUse = React.lazy(() => import('./model-use/apriori/model-use'));
  const CosineModelUse = React.lazy(() => import('./model-use/cosine/cosine-model-use'));
  const pageTitle = useAppSelector(state => state.app.pageTitle);
  const selectedMenuKey = useAppSelector(state => state.app.selectedMenuKey);
  const selectedModelType = useAppSelector(s => s.app.selectedModelType);
  const dispatch = useAppDispatch();
  const { Title } = Typography;

  function handleMenuChange(key: MenuKey) {
    dispatch(changeSelectedMenuKey(key));
  }

  return (
    <>
      <Router basename="/app">
        <Layout className="layout">
          <Header>
            <div className="logo" />
            <Menu theme="dark" mode="horizontal" selectedKeys={[selectedMenuKey]} onSelect={e => handleMenuChange(e.key as MenuKey)}>
              <Menu.Item key={MenuKey.DataImport}><Link to="/data-import">Data Import</Link></Menu.Item>
              <Menu.Item key={MenuKey.ModelStatistics}><Link to={`/model-stats/${selectedModelType}`}>Model Statistics</Link></Menu.Item>
              <Menu.Item key={MenuKey.ModelUsage}><Link to={`/model-use/${selectedModelType}`}>Model Usage</Link></Menu.Item>
            </Menu>
          </Header>
          <Content style={{ padding: '0 50px' }}>
            <div className="page-title-container">
              <Title level={3}>{pageTitle}</Title>
            </div>
            <div className="site-layout-content">
              <Suspense fallback={<Loading />}>
                <div>
                  <Switch>
                    <Route exact path="/">
                      <Redirect to="/data-import" />
                    </Route>
                    <Route path="/data-import">
                      <DataImport />
                    </Route>
                    <Route path="/model-stats/apriori">
                      <ModelStats />
                    </Route>
                    <Route path="/model-stats/cosine">
                      <CosineModelStats />
                    </Route>
                    <Route path="/model-use/apriori">
                      <ModelUse />
                    </Route>
                    <Route path="/model-use/cosine">
                      <CosineModelUse />
                    </Route>
                  </Switch>
                </div>
              </Suspense>
            </div>
          </Content>
          <Footer style={{ textAlign: 'center' }}>Created by applied-math-coding</Footer>
        </Layout>
      </Router>
    </>
  );
}

export default App;
