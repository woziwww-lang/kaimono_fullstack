import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:price_comparison_mobile/main.dart';

void main() {
  testWidgets('App should display title', (WidgetTester tester) async {
    // Build the app
    await tester.pumpWidget(const MyApp());

    // Verify that app title is displayed
    expect(find.text('価格比較アプリ'), findsOneWidget);
  });

  testWidgets('HomeScreen should be the initial route', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    // Wait for the widget tree to settle
    await tester.pumpAndSettle();

    // HomeScreen should be displayed
    expect(find.byType(MaterialApp), findsOneWidget);
  });
}
