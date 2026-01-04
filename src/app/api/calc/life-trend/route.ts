/**
 * Life Trend Calculation API
 * 人生趋势计算接口
 */

import { NextRequest, NextResponse } from 'next/server';
import { 
  calculateNatalChart, 
  calculateLifeTrend,
  type BirthData 
} from '@/lib/astro';

/**
 * POST /api/calc/life-trend
 */
export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    
    const { birthData: bd, startYear, endYear } = body;
    
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
    
    const birthYear = birthData.date.getFullYear();
    const start = startYear || birthYear;
    const end = endYear || birthYear + 80;
    
    const trend = calculateLifeTrend(chart, start, end, 'yearly');
    
    return NextResponse.json({
      type: 'yearly',
      birthDate: birthData.date.toISOString(),
      startYear: start,
      endYear: end,
      points: trend.points.map(p => ({
        date: p.date.toISOString(),
        year: p.year,
        age: p.age,
        overallScore: p.overallScore,
        harmonious: p.harmonious,
        challenge: p.challenge,
        transformation: p.transformation,
        dimensions: p.dimensions,
        dominantPlanet: p.dominantPlanet,
        profection: p.profection,
        isMajorTransit: p.isMajorTransit,
        majorTransitName: p.majorTransitName,
        lunarPhaseName: p.lunarPhaseName,
      })),
      summary: trend.summary,
      cycles: {
        saturnCycles: trend.cycles.saturnCycles,
        jupiterCycles: trend.cycles.jupiterCycles,
        profectionCycles: trend.cycles.profectionCycles,
      },
    });
  } catch (error) {
    console.error('Life trend error:', error);
    return NextResponse.json(
      { error: 'Failed to calculate life trend' },
      { status: 500 }
    );
  }
}

