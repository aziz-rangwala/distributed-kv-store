import React, { useState } from 'react';
import KeyValueForm from './components/KeyValueForm';
import SearchKey from './components/SearchKey';

export default function App() {
  return (
    <div className="app">
      <h1>Distributed Key-Value Store</h1>
      <KeyValueForm />
      <SearchKey />
    </div>
  );
}