/**
 * Weekly Forecast Calculation API
 * 每周预测计算接口
 */

import { NextRequest, NextResponse } from 'next/server';
import { 
  calculateNatalChart, 
  processWeeklyForecast,
  calculateWeeklyForecast,
  type BirthData 
} from '@/lib/astro';

/**
 * POST /api/calc/weekly
 */
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const { birthData: bd, date, withFactors = true } = body;
    
    if (!bd || !bd.date || bd.latitude === undefined || bd.longitude === undefined) {
      return NextResponse.json(
        { error: 'Missing required birthData fields' },
        { status: 400 }
      );
    }
    
    const birthData: BirthData = {
      name: bd.name || 'Anonymous',
      date: new Date(bd.date),
      latitude: bd.latitude,
      longitude: bd.longitude,
      timezone: bd.timezone || 'UTC',
    };
    
    const chart = calculateNatalChart(birthData);
    
    // 获取本周一
    const targetDate = date ? new Date(date) : new Date();
    const dayOfWeek = targetDate.getDay();
    const diff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
    const weekStart = new Date(targetDate);
    weekStart.setDate(weekStart.getDate() + diff);
    
    if (withFactors) {
      const forecast = processWeeklyForecast(chart, weekStart);
      
      return NextResponse.json({
        type: 'weekly',
        processed: true,
        startDate: forecast.startDate.toISOString(),
        endDate: forecast.endDate.toISOString(),
        overallScore: forecast.overallScore,
        weeklyTheme: forecast.weeklyTheme,
        weeklyInsight: forecast.weeklyInsight,
        weeklyFactors: forecast.weeklyFactors,
        dimensionTrends: forecast.dimensionTrends,
        days: forecast.days.map(d => ({
          date: d.date.toISOString(),
          dayOfWeek: d.dayOfWeek,
          overallScore: d.overallScore,
          rawScore: d.rawScore,
          dimensions: d.dimensions,
          topFactors: d.topFactors?.slice(0, 3),
        })),
      });
    } else {
      const forecast = calculateWeeklyForecast(weekStart, chart);
      
      return NextResponse.json({
        type: 'weekly',
        processed: false,
        startDate: forecast.startDate.toISOString(),
        endDate: forecast.endDate.toISOString(),
        overallTheme: forecast.overallTheme,
        keyDates: forecast.keyDates,
        dailySummaries: forecast.dailySummaries,
      });
    }
  } catch (error) {
    console.error('Weekly forecast error:', error);
    return NextResponse.json(
      { error: 'Failed to calculate weekly forecast' },
      { status: 500 }
    );
  }
}

