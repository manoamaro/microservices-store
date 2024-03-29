import * as React from 'react';
import Head from 'next/head';
import {AppProps} from 'next/app';
import {ThemeProvider} from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import {CacheProvider, EmotionCache} from '@emotion/react';
import theme from '../src/theme';
import createEmotionCache from '../src/createEmotionCache';
import RouteGuard from "../src/RouteGuard";
import {Box} from "@mui/material";
import Bar from "../src/Bar";
import {AuthProvider} from "../src/AuthProvider";

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

export interface MyAppProps extends AppProps {
    emotionCache?: EmotionCache;
}

export default function MyApp(props: MyAppProps) {
    const {Component, emotionCache = clientSideEmotionCache, pageProps} = props;

  return (
      <AuthProvider>
          <CacheProvider value={emotionCache}>
              <Head>
                  <meta name="viewport" content="initial-scale=1, width=device-width"/>
              </Head>
              <ThemeProvider theme={theme}>
                  <CssBaseline/>
                  <Box sx={{flexGrow: 1}}>
                      <Bar/>
                      <RouteGuard>
                          <Component {...pageProps} />
                      </RouteGuard>
                  </Box>
              </ThemeProvider>
          </CacheProvider>
      </AuthProvider>
  );
}
