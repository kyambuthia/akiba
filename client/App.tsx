import React from 'react';
import { SafeAreaView, ScrollView, StyleSheet, Text, View } from 'react-native';
import { SignupScreen } from './src/screens/SignupScreen';
import { LoginScreen } from './src/screens/LoginScreen';

export default function App() {
  return (
    <SafeAreaView style={styles.safe}>
      <ScrollView contentContainerStyle={styles.container}>
        <Text style={styles.title}>Akiba Client Scaffold</Text>
        <View style={styles.card}><SignupScreen /></View>
        <View style={styles.card}><LoginScreen /></View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safe: { flex: 1, backgroundColor: '#f7f4ee' },
  container: { padding: 20, gap: 14 },
  title: { fontSize: 24, fontWeight: '700', color: '#273043' },
  card: { backgroundColor: '#fff', borderRadius: 10, padding: 14, borderWidth: 1, borderColor: '#ddd' }
});
